package services

import "time"

type CacheService interface {
	GetObject(key string, dest any) (bool, error)
	GetValue(key string) (string, bool, error)

	SetObject(key string, obj any, exp time.Duration) error
	SetValue(key string, value string, exp time.Duration) error

	RemoveKey(ck string) error
	RemoveKeyWithCount(ck string) (int64, error)
	RemoveKeysWithCount(cks []string) (int64, error)

	AddSet(setKey string, member string) error
	GetSetMembers(setKey string) ([]string, error)
	RemoveSetMember(setKey string, member string) error
	RemoveKeys(keys ...string) error

	SetH(key string, v map[string]any, exp time.Duration) error
	GetH(key string, field string) (string, error)
}

type HasIsActiveStatus interface {
	GetIsActive() bool
}

type FormErrors map[string]error
