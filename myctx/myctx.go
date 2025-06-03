package myctx

import (
	"context"
	"errors"
)

type myContextKey string

const (
	userIdKey myContextKey = "userId"
	shopIdKey myContextKey = "shopId"
	token     myContextKey = "token"
	isAuth    myContextKey = "isAuth"
)

func GetIdsFromContext(ctx context.Context) (int, string, error) {
	userId, ok := ctx.Value(userIdKey).(int)
	if !ok {
		return 0, "", errors.New("userId is required")
	}
	shopId, ok := ctx.Value(shopIdKey).(string)
	if !ok {
		return 0, "", errors.New("shopId is required")
	}

	return userId, shopId, nil
}

func SetIds(ctx context.Context, userId int, shopId string) context.Context {
	ctx = context.WithValue(ctx, userIdKey, userId)
	ctx = context.WithValue(ctx, shopIdKey, shopId)
	return ctx
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
