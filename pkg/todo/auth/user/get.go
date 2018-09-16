package user

import (
	gocontext "context"
	"fmt"

	"github.com/golangci/golangci-api/app/models"
	"github.com/golangci/golangci-api/pkg/todo/auth/sess"
	"github.com/golangci/golangci-api/pkg/todo/db"
	"github.com/golangci/golib/server/context"
	"github.com/golangci/golib/server/handlers/herrors"

	"github.com/jinzhu/gorm"
)

type userCtxKeyType string

var userCtxKey userCtxKeyType = "user"
var githubAuthCtxKey userCtxKeyType = "githubAuth"

var ErrNotAuthorized = herrors.New403Errorf("user isn't authorized")

func GetCurrentID(ctx *context.C) (uint, error) {
	if ctx.R == nil { // background, no request
		return 0, fmt.Errorf("no user for background processing")
	}

	userIDi, err := sess.GetValue(ctx, userIDSessKey)
	if err != nil {
		return 0, err
	}

	if userIDi == nil {
		return 0, ErrNotAuthorized
	}

	userIDf := userIDi.(float64)
	return uint(userIDf), nil
}

func GetCurrent(ctx *context.C) (*models.User, error) {
	if v := ctx.Ctx.Value(userCtxKey); v != nil {
		user := v.(models.User)
		return &user, nil
	}

	userID, err := GetCurrentID(ctx)
	if err != nil {
		return nil, err
	}

	var u models.User
	err = models.NewUserQuerySet(db.Get(ctx)).IDEq(userID).One(&u)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// deleted test user
			return nil, herrors.New403Errorf("can't get current user with id %d: %s", userID, err)
		}

		return nil, fmt.Errorf("can't get current user with id %d: %s", userID, err)
	}

	ctx.Ctx = gocontext.WithValue(ctx.Ctx, userCtxKey, u)
	return &u, nil
}

func GetGithubAuth(ctx *context.C) (*models.GithubAuth, error) {
	if v := ctx.Ctx.Value(githubAuthCtxKey); v != nil {
		a := v.(*models.GithubAuth)
		return a, nil
	}

	userID, err := GetCurrentID(ctx)
	if err != nil {
		return nil, err
	}

	ga, err := GetGithubAuthForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	ctx.Ctx = gocontext.WithValue(ctx.Ctx, githubAuthCtxKey, ga)
	return ga, nil
}

func GetGithubAuthForUser(ctx *context.C, userID uint) (*models.GithubAuth, error) {
	var ga models.GithubAuth
	err := models.NewGithubAuthQuerySet(db.Get(ctx)).
		UserIDEq(userID).
		OrderDescByID().
		One(&ga)
	if err != nil {
		return nil, herrors.New(err, "can't get github auth for user %d", userID)
	}

	return &ga, nil
}