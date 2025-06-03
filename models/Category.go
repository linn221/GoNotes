package models

import (
	"gorm.io/gorm"
)

type Category struct {
	Id          int     `gorm:"primaryKey"`
	Name        string  `gorm:"index;not null"`
	Description *string `gorm:"default:null"`
	HasShopId
	HasIsActive
}

type CategoryResource struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	HasShopId
}

type CategoryDetailResource struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	HasIsActive
	HasShopId
}

// func (cat *Category) AfterCreate(db *gorm.DB) error {
// 	if err := ; err != nil {
// 		return err
// 	}

// clean related cache when category is updated
func (readers *ReadServices) AfterCategoryUpdate(db *gorm.DB, shopId string, id int) error {

	if err := readers.CategoryGetService.CleanCache(id); err != nil {
		return err
	}
	if err := readers.CategoryListService.CleanCache(shopId); err != nil {
		return err
	}
	var relatedItemIds []int
	if err := db.Model(&Item{}).Where("category_id = ?", id).Pluck("id", &relatedItemIds).Error; err != nil {
		return err
	}
	for _, itemId := range relatedItemIds {
		if err := readers.ItemGetService.CleanCache(itemId); err != nil {
			return err
		}
	}
	if err := readers.ItemListService.CleanCache(shopId); err != nil {
		return err
	}
	return nil

}
