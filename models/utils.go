package models

import "gorm.io/gorm"

func first[T any](db *gorm.DB, userId int, id int, preloads ...string) (*T, error) {
	var v T
	dbCtx := db.Where("user_id = ?", userId)
	if len(preloads) > 0 {
		dbCtx = dbCtx.Preload(preloads[0])
		for _, p := range preloads[1:] {
			dbCtx.Preload(p)
		}
	}

	if err := dbCtx.First(&v, id).Error; err != nil {
		return nil, err
	}

	return &v, nil
}

func find[T any](db *gorm.DB, userId int, preloads ...string) ([]T, error) {
	var results []T
	dbCtx := db.Where("user_id = ?", userId)
	if len(preloads) > 0 {
		dbCtx = dbCtx.Preload(preloads[0])
		for _, p := range preloads[1:] {
			dbCtx.Preload(p)
		}
	}

	if err := dbCtx.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
