package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func NewSession(userId int, shopId string, cache CacheService) (string, error) {
	token := uuid.NewString()
	cacheVal := fmt.Sprintf("%d:%s", userId, shopId)
	cacheKey := "Token:" + token
	err := cache.SetValue(cacheKey, cacheVal, time.Hour*127)
	return token, err
}

func RemoveSession(token string, cache CacheService) error {
	return cache.RemoveKey("Token:" + token)
}
