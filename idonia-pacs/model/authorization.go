package model

import (
	"context"
	"github.com/pkg/errors"
)

type Authorization struct {
	Token string
}

type LoginRQ struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRS struct {
	Token string `json:"token"`
}

func GetAuthorizationFromContext(ctx context.Context) (auth *Authorization, err error) {
	if auth, ok := ctx.Value("authorization").(*Authorization); ok {
		return auth, nil
	}
	return nil, errors.New("No authorization object found")
}

const PermissionLevelAll = "all"
