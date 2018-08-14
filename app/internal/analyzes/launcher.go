package analyzes

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/satori/go.uuid"

	"github.com/golangci/golangci-api/app/internal/db"
	"github.com/golangci/golangci-api/app/internal/errors"
	"github.com/golangci/golangci-api/app/models"
	"github.com/golangci/golangci-api/app/utils"
	"github.com/golangci/golangci-worker/app/analyze/analyzequeue"
	"github.com/golangci/golangci-worker/app/analyze/analyzequeue/task"
	"github.com/golangci/golib/server/context"
)

// reanalyze each repo every reanalyzeInterval duration
var reanalyzeInterval = getDurationFromEnv("REPO_REANALYZE_INTERVAL", 30*time.Minute)

func getDurationFromEnv(key string, def time.Duration) time.Duration {
	cfgStr := os.Getenv(key)
	if cfgStr == "" {
		return def
	}

	d, err := time.ParseDuration(cfgStr)
	if err != nil {
		logrus.Errorf("Invalid %s %q: %s", key, cfgStr, err)
		return def
	}

	return d
}

func StartLauncher() {
	go launchRepoAnalyzes()
	go restartRepoAnalyzes()
}

func launchRepoAnalyzes() {
	ctx := utils.NewBackgroundContext()

	checkInterval := reanalyzeInterval / 2
	const minCheckInterval = time.Minute * 5
	if checkInterval < minCheckInterval {
		checkInterval = minCheckInterval
	}

	for range time.Tick(checkInterval) {
		if err := launchAnalyzes(ctx); err != nil {
			errors.Warnf(ctx, "Can't launch analyzes: %s", err)
			continue
		}
	}
}

func createNewAnalysisStatuses(ctx *context.C) error {
	var analysisStatuses []models.RepoAnalysisStatus
	err := models.NewRepoAnalysisStatusQuerySet(db.Get(ctx)).
		All(&analysisStatuses)
	if err != nil {
		return fmt.Errorf("can't get all analysis statuses: %s", err)
	}

	repoToStatus := map[string]models.RepoAnalysisStatus{}
	for _, as := range analysisStatuses {
		repoToStatus[strings.ToLower(as.Name)] = as
	}

	var githubRepos []models.GithubRepo
	err = models.NewGithubRepoQuerySet(db.Get(ctx)).
		All(&githubRepos)
	if err != nil {
		return fmt.Errorf("can't get all github repos: %s", err)
	}

	for _, repo := range githubRepos {
		_, ok := repoToStatus[strings.ToLower(repo.Name)]
		if ok {
			continue
		}

		as := models.RepoAnalysisStatus{
			Name: strings.ToLower(repo.Name),
		}
		if err = as.Create(db.Get(ctx)); err != nil {
			return fmt.Errorf("can't create repo analysis status %+v: %s", as, err)
		}
	}

	return nil
}

func launchAnalyzes(ctx *context.C) error {
	if err := createNewAnalysisStatuses(ctx); err != nil {
		return fmt.Errorf("can't create new analysis statuses: %s", err)
	}

	var analysisStatuses []models.RepoAnalysisStatus
	err := models.NewRepoAnalysisStatusQuerySet(db.Get(ctx)).
		All(&analysisStatuses)
	if err != nil {
		return fmt.Errorf("can't get all analysis statuses: %s", err)
	}

	for _, as := range analysisStatuses {
		if err := processAnalysisStatus(ctx, &as); err != nil {
			return err
		}
	}

	return nil
}

func processAnalysisStatus(ctx *context.C, as *models.RepoAnalysisStatus) error {
	needAnalysis := as.LastAnalyzedAt.IsZero() ||
		(as.HasPendingChanges && as.LastAnalyzedAt.Add(reanalyzeInterval).Before(time.Now()))
	if !needAnalysis {
		ctx.L.Infof("No need to launch analysis for analysis status %v: last_analyzed=%s ago, reanalyze_interval=%s",
			as, time.Since(as.LastAnalyzedAt), reanalyzeInterval)
		return nil
	}

	ctx.L.Infof("Launching analysis for %+v...", as)
	if err := launchAnalysis(ctx, as); err != nil {
		return fmt.Errorf("can't launch analysis %+v: %s", as, err)
	}

	return nil
}

func OnRepoMasterUpdated(ctx *context.C, repoName, defaultBranch, commitSHA string) error {
	var as models.RepoAnalysisStatus
	err := models.NewRepoAnalysisStatusQuerySet(db.Get(ctx)).
		NameEq(repoName).
		One(&as)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			as = models.RepoAnalysisStatus{
				Name: repoName,
			}
			if err = as.Create(db.Get(ctx)); err != nil {
				return fmt.Errorf("can't create repo analysis status %+v: %s", as, err)
			}
		} else {
			return fmt.Errorf("can't fetch analysis status with name %s: %s", repoName, err)
		}
	}

	as.HasPendingChanges = true
	as.DefaultBranch = defaultBranch
	as.PendingCommitSHA = commitSHA
	err = as.Update(db.Get(ctx),
		models.RepoAnalysisStatusDBSchema.HasPendingChanges,
		models.RepoAnalysisStatusDBSchema.DefaultBranch,
		models.RepoAnalysisStatusDBSchema.PendingCommitSHA,
	)
	if err != nil {
		return fmt.Errorf("can't update has_pending_changes to true: %s", err)
	}

	ctx.L.Infof("Set has_pending_changes=true for repo %s analysis status", repoName)

	return processAnalysisStatus(ctx, &as)
}

func launchAnalysis(ctx *context.C, as *models.RepoAnalysisStatus) (err error) {
	finishTx, err := db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer finishTx(&err)

	a := models.RepoAnalysis{
		RepoAnalysisStatusID: as.ID,
		AnalysisGUID:         uuid.NewV4().String(),
		Status:               "sent_to_queue",
		CommitSHA:            as.PendingCommitSHA,
		ResultJSON:           []byte("{}"),
	}
	if err = a.Create(db.Get(ctx)); err != nil {
		return fmt.Errorf("can't create repo analysis: %s", err)
	}

	t := &task.RepoAnalysis{
		Name:         strings.ToLower(as.Name),
		AnalysisGUID: a.AnalysisGUID,
		Branch:       as.DefaultBranch,
	}

	if err = analyzequeue.ScheduleRepoAnalysis(t); err != nil {
		return fmt.Errorf("can't send repo for analysis into queue: %s", err)
	}

	n, err := models.NewRepoAnalysisStatusQuerySet(db.Get(ctx)).
		NameEq(strings.ToLower(as.Name)).
		VersionEq(as.Version).
		GetUpdater().
		SetHasPendingChanges(false).
		SetPendingCommitSHA("").
		SetVersion(as.Version + 1).
		SetLastAnalyzedAt(time.Now().UTC()).
		UpdateNum()
	if err != nil {
		return fmt.Errorf("can't update repo analysis status after processing: %s", err)
	}
	if n == 0 {
		return fmt.Errorf("got race condition updating repo analysis status on version %d->%d",
			as.Version, as.Version+1)
	}

	return nil
}

func restartRepoAnalyzes() {
	repoAnalysisTimeout := getDurationFromEnv("REPO_ANALYSIS_TIMEOUT", 15*time.Minute)
	ctx := utils.NewBackgroundContext()

	for range time.Tick(repoAnalysisTimeout / 2) {
		if err := runRestartRepoAnalyzesIter(ctx, repoAnalysisTimeout); err != nil {
			errors.Warnf(ctx, "Can't restart analyzes: %s", err)
		}
	}
}

func runRestartRepoAnalyzesIter(ctx *context.C, repoAnalysisTimeout time.Duration) error {
	var analyzes []models.RepoAnalysis
	err := models.NewRepoAnalysisQuerySet(db.Get(ctx)).
		StatusIn("sent_to_queue", "processing", "error").
		CreatedAtLt(time.Now().Add(-repoAnalysisTimeout)).
		PreloadRepoAnalysisStatus().
		All(&analyzes)
	if err != nil {
		return fmt.Errorf("can't get repo analyzes: %s", err)
	}

	if len(analyzes) == 0 {
		return nil
	}

	for _, a := range analyzes {
		as := a.RepoAnalysisStatus
		t := &task.RepoAnalysis{
			Name:         strings.ToLower(as.Name),
			AnalysisGUID: a.AnalysisGUID,
			Branch:       as.DefaultBranch,
		}

		if err = analyzequeue.ScheduleRepoAnalysis(t); err != nil {
			return fmt.Errorf("can't resend repo %s for analysis into queue: %s", as.Name, err)
		}

		errors.Warnf(ctx, "Restarted analysis for %s in status %s", as.Name, a.Status)
	}

	return nil
}
