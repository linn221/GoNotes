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
}

type Getter[T any] interface {
	Get(shopId string, id int) (*T, bool, error)
	CleanCache(id int) error
}
type Lister[T any] interface {
	List(shopId string) ([]T, error)
	CleanCache(shopId string) error
}

type CleanInstanceCache func(id int) error
type CleanListingCache func(shopId string) error

type HasIsActiveStatus interface {
	GetIsActive() bool
}
