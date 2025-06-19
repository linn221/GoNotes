package myctx

import (
	"context"
	"errors"
)

type myContextKey string

const (
	userIdKey myContextKey = "userId"
	token     myContextKey = "token"
	isAuth    myContextKey = "isAuth"
)

func SetUserId(ctx context.Context, userId int) context.Context {
	ctx = context.WithValue(ctx, userIdKey, userId)
	return ctx
}

func GetUserId(ctx context.Context) (int, error) {
	userId, ok := ctx.Value(userIdKey).(int)
	if !ok || userId < 1 {
		return 0, errors.New("userId is required")
	}
	return userId, nil
}

func IsAuth(ctx context.Context) bool {
	isAuth, ok := ctx.Value(isAuth).(bool)
	if !ok || !isAuth {
		return false
	}
	return true
}

func SetAuth(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, isAuth, true)
	return ctx
}
