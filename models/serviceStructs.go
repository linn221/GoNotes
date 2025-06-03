package models

import (
	"fmt"
	"linn221/shop/services"
	"time"

	"gorm.io/gorm"
)

/*
	generic GetService for internal use
	retrieve from redis or fetch from db (& store in redis for later)
*/

const forever = time.Duration(0)

type Resource interface {
	GetShopId() string
}

type defaultGetService[T Resource] struct {
	db          *gorm.DB
	cache       services.CacheService
	table       string
	cachePrefix string
	cacheLength time.Duration
}

type CacheValue[T any] struct { // a workaround as T type exclude ShopId field from json
	Object T
	ShopId string
}

func (service *defaultGetService[T]) Get(shopId string, id int) (*T, bool, error) {
	var cacheValue CacheValue[T]

	redisKey := service.cachePrefix + ":" + fmt.Sprint(id)
	exists, err := service.cache.GetObject(redisKey, &cacheValue)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		var result T
		if err := service.db.Table(service.table).
			First(&result, id).Error; err != nil { // using gorm's smart scan
			if err == gorm.ErrRecordNotFound {
				return nil, false, nil
			}
			return nil, false, err
		}
		shopID := result.GetShopId()
		cacheValue = CacheValue[T]{Object: result, ShopId: shopID}
		if err := service.cache.SetObject(redisKey, &cacheValue, service.cacheLength); err != nil {
			return nil, false, err
		}
	}
	if cacheValue.ShopId != shopId { // don't let user access resource owned by another shop
		return nil, false, nil
	}

	return &cacheValue.Object, true, nil
}

// clear cache
func (service *defaultGetService[T]) CleanCache(id int) error {
	key := service.cachePrefix + ":" + fmt.Sprint(id)
	return service.cache.RemoveKey(key)
}

// only listing active ones
type defaultListService[T any] struct {
	db          *gorm.DB
	cache       services.CacheService
	table       string
	cachePrefix string
	cacheLength time.Duration
}

func (service *defaultListService[T]) List(shopId string) ([]T, error) {
	var results []T

	key := service.cachePrefix + ":" + shopId
	exists, err := service.cache.GetObject(key, &results)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := service.db.Table(service.table).Where("shop_id = ? AND is_active = 1", shopId).Find(&results).Error; err != nil {
			return nil, err
		}
		if err := service.cache.SetObject(key, &results, service.cacheLength); err != nil {
			return nil, err
		}
	}

	return results, nil
}

// clear cache
func (service *defaultListService[ResponseT]) CleanCache(shopId string) error {
	return service.cache.RemoveKey(service.cachePrefix + ":" + fmt.Sprint(shopId))
}

type customGetService[T Resource] struct {
	cachePrefix string
	cacheLength time.Duration
	db          *gorm.DB
	cache       services.CacheService
	fetch       func(db *gorm.DB, id int) (T, error)
}

func (service *customGetService[T]) Get(shopId string, id int) (*T, bool, error) {
	var cacheValue CacheValue[T]

	redisKey := service.cachePrefix + ":" + fmt.Sprint(id)
	exists, err := service.cache.GetObject(redisKey, &cacheValue)
	if err != nil {
		return nil, false, err
	}

	if !exists {
		result, err := service.fetch(service.db, id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, false, nil
			}
			return nil, false, err
		}
		cacheValue = CacheValue[T]{Object: result, ShopId: result.GetShopId()}
		if err := service.cache.SetObject(redisKey, &cacheValue, service.cacheLength); err != nil {
			return nil, false, err
		}
	}

	if cacheValue.ShopId != shopId {
		return nil, false, nil // only saying not found if belongs to another shop
	}

	return &cacheValue.Object, true, nil
}
func (service *customGetService[T]) CleanCache(id int) error {
	return service.cache.RemoveKey(service.cachePrefix + ":" + fmt.Sprint(id))
}

// only listing active ones
type customListService[T any] struct {
	db          *gorm.DB
	cache       services.CacheService
	cachePrefix string
	cacheLength time.Duration
	fetch       func(db *gorm.DB, shopId string) ([]T, error)
}

func (service *customListService[T]) List(shopId string) ([]T, error) {
	var results []T

	key := service.cachePrefix + ":" + shopId
	exists, err := service.cache.GetObject(key, &results)
	if err != nil {
		return nil, err
	}
	if !exists {
		results, err = service.fetch(service.db, shopId)
		if err != nil {
			return nil, err
		}
		if err := service.cache.SetObject(key, &results, service.cacheLength); err != nil {
			return nil, err
		}
	}

	return results, nil
}

// clear cache
func (service *customListService[ResponseT]) CleanCache(shopId string) error {
	return service.cache.RemoveKey(service.cachePrefix + ":" + shopId)
}

// type customListService[T any] struct {
// 	cacheKey    config.CacheKey
// 	cacheLength time.Duration
// 	fetch       func(db *gorm.DB) ([]T, error)
// }

// // may inject queries by db
// func (service *customListService[ResponseT]) List(db *gorm.DB) ([]ResponseT, *myerror.ServiceError) {
// 	var results []ResponseT

// 	exists, err := config.GetRedisObject(service.cacheKey, &results)
// 	if err != nil {
// 		return results, myerror.NewWithMessage("cache error", err)
// 	}
// 	if !exists {
// 		results, err = service.fetch(db)
// 		if err != nil {
// 			if err == gorm.ErrRecordNotFound {
// 				return results, &myerror.ServiceError{
// 					StatusCode: http.StatusNotFound,
// 					Message:    "record not found",
// 				}
// 			}
// 			return results, myerror.NewWithMessage("error fetching database", err)
// 		}
// 		if err := config.SetRedisObject(service.cacheKey, results, service.cacheLength); err != nil {
// 			return results, myerror.NewWithMessage("cache error", err)
// 		}
// 	}

// 	return results, nil
// }

// func (service *customListService[T]) clear() error {
// 	return config.RemoveRedisKey(service.cacheKey)
// }

// type paginateSerivce[T any] struct {
// 	cachePrefix cachePrefix
// 	limit       int
// 	// order       string
// 	cacheLength time.Duration
// 	fetch       func(db *gorm.DB, page int) ([]T, error)
// }

// type PaginateResponse[T any] struct {
// 	Page        int
// 	Records     []T
// 	HasNextPage bool
// }

// func (service *paginateSerivce[T]) Paginate(db *gorm.DB, page int) (*PaginateResponse[T], error) {
// 	records, err := service.fetch(db, page)
// 	if err != nil {
// 		// if err == gorm.ErrRecordNotFound {
// 		// 	return nil, &myerror.ServiceError{
// 		// 		StatusCode: http.StatusNotFound,
// 		// 		Message:    "record not found",
// 		// 	}
// 		// }
// 		// return results, myerror.NewWithMessage("error fetching database", err)
// 		return nil, err
// 	}

// 	var results PaginateResponse[T]
// 	if len(records) == service.limit {
// 		records = records[:len(records)-1]
// 		results.HasNextPage = true
// 	}
// 	results.Page = page
// 	results.Records = records

// 	return &results, nil
// }

// func (service *paginateSerivce[T]) clear() error {
// 	return RemoveSequentialKeys(service.cachePrefix)
// }

// func RemoveSequentialKeys(cp cachePrefix) error {
// 	i := 1
// 	for {
// 		count, err := config.RemoveRedisKeyWithCount(cp.cacheKey(i))
// 		if err != nil {
// 			return err
// 		}
// 		if count == 0 {
// 			break
// 		}
// 		i++
// 	}

// 	return nil
// }

// // func (service *paginateSerivce[T]) PaginateWithCache(db *gorm.DB, page int) ([]PaginateResponse[T], error) {

// // 	return nil, nil
// // }

// type HasKey[K comparable] interface {
// 	InternalMapKey() K
// }

// type internalMap[K comparable, T HasKey[K]] struct {
// 	cacheKey    config.CacheKey
// 	cacheLength time.Duration
// 	fetch       func(db *gorm.DB) ([]T, error)
// 	// key         func(T) K
// }

// func (service *internalMap[K, T]) GetMap(db *gorm.DB) (map[K]T, error) {
// 	var m map[K]T
// 	found, err := config.GetRedisObject(service.cacheKey, &m)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !found {
// 		records, err := service.fetch(db)
// 		if err != nil {
// 			return nil, err
// 		}
// 		m = make(map[K]T, len(records))
// 		for _, record := range records {
// 			m[record.InternalMapKey()] = record
// 		}
// 		if err := config.SetRedisObject(service.cacheKey, m, service.cacheLength); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return m, nil
// }

// func (service *internalMap[K, T]) clear() error {
// 	return config.RemoveRedisKey(service.cacheKey)
// }

// type Transformer[T any] interface {
// 	Transform(*gorm.DB) (T, error)
// }

// type SearchService[S fmt.Stringer, T any, C Transformer[T]] struct {
// 	Key         config.CacheKey
// 	CachePrefix cachePrefixEI
// 	Length      time.Duration
// 	Limit       int
// 	Fetch       func(*gorm.DB, S) (C, error)
// }

// func (service SearchService[S, T, C]) Search(db *gorm.DB, params S) (T, *myerror.ServiceError) {

// 	var v T
// 	hash := utils.HashString(params.String())
// 	var cache C
// 	key := service.CachePrefix.cacheKey(hash)
// 	ok, err := config.GetRedisObject(key, &cache)
// 	if err != nil {
// 		return v, myerror.NewWithMessage("error", err)
// 	}
// 	if !ok {
// 		cache, err = service.Fetch(db, params)
// 		if err != nil {
// 			return v, myerror.NewWithMessage("error fetching database", err)
// 		}
// 		if err := config.SetRedisObject(key, &cache, service.Length); err != nil {
// 			return v, myerror.NewWithMessage("redis error", err)
// 		}
// 		if err := config.AddRedisSet(service.Key, hash); err != nil {
// 			return v, myerror.NewWithMessage("redis eror", err)
// 		}
// 	}

// 	results, err := cache.Transform(db)
// 	if err != nil {
// 		return v, myerror.NewWithMessage("redis error", err)
// 	}

// 	return results, nil
// }

// func (service SearchService[S, T, C]) clear() error {
// 	searchResultHashes, err := config.GetRedisSetMembers(service.Key)
// 	if err != nil {
// 		return err
// 	}
// 	if len(searchResultHashes) > 0 {

// 		searchResultKeys := make([]string, 0, len(searchResultHashes))
// 		for _, key := range searchResultHashes {
// 			searchResultKeys = append(searchResultKeys, string(service.CachePrefix.cacheKey(key)))
// 		}

// 		count, err := config.RemoveRedisKeysWithCount(searchResultKeys)
// 		if err != nil {
// 			return err
// 		}
// 		if count != int64(len(searchResultKeys)) {
// 			return errors.New("error deleting some of hashes")
// 		}
// 	}
// 	if err := config.RemoveRedisKey(service.Key); err != nil {
// 		return err
// 	}

// 	return nil
// }
