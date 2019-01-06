package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

//go:generate goqueryset -in repo.go

type UniversalRepo interface {
	Owner() string
	Repo() string
}

type RepoCommitState string

const (
	RepoCommitStateCreateInit        RepoCommitState = "create/init"
	RepoCommitStateCreateSentToQueue RepoCommitState = "create/sent_to_queue"
	RepoCommitStateCreateCreatedRepo RepoCommitState = "create/created_repo"
	RepoCommitStateCreateDone        RepoCommitState = "create/done"

	RepoCommitStateDeleteInit        RepoCommitState = "delete/init"
	RepoCommitStateDeleteSentToQueue RepoCommitState = "delete/sent_to_queue"
	RepoCommitStateDeleteDone        RepoCommitState = "delete/done"
)

func (s RepoCommitState) IsDeleteState() bool {
	return s == RepoCommitStateDeleteInit || s == RepoCommitStateDeleteSentToQueue || s == RepoCommitStateDeleteDone
}

func (s RepoCommitState) IsCreateState() bool {
	return s == RepoCommitStateCreateInit || s == RepoCommitStateCreateSentToQueue ||
		s == RepoCommitStateCreateCreatedRepo || s == RepoCommitStateCreateDone
}

func (s RepoCommitState) IsDone() bool {
	return s == RepoCommitStateDeleteDone || s == RepoCommitStateCreateDone
}

//gen:qs
type Repo struct {
	gorm.Model

	// TODO: remove it and move repo connections to another table,
	// take organizations into account
	UserID uint // user who the last time connected this repo

	FullName        string `gorm:"column:name"`         // lower-cased DisplayName to avoid case-sensitivity bugs
	DisplayFullName string `gorm:"column:display_name"` // original name of repo from github: original register is saved

	HookID string

	Provider       string // github.com, gitlab.com etc
	ProviderHookID int
	ProviderID     int // provider repo id: use it (not name) as repo identifier because of repo renaming

	CommitState RepoCommitState // state of creation or deletion

	StargazersCount int
	IsPrivate       bool
}

func (r *Repo) Owner() string {
	return strings.ToLower(strings.Split(r.FullName, "/")[0])
}

func (r *Repo) Repo() string {
	return strings.ToLower(strings.Split(r.FullName, "/")[1])
}

func (r *Repo) String() string {
	return r.FullName
}

func (r *Repo) GoString() string {
	return fmt.Sprintf("{Name: %s, ID: %d, CommitState: %s}", r.FullName, r.ID, r.CommitState)
}

func (r Repo) IsDeleting() bool {
	return r.CommitState.IsDeleteState() && !r.CommitState.IsDone()
}

func (r Repo) IsCreating() bool {
	return r.CommitState.IsCreateState() && !r.CommitState.IsDone()
}