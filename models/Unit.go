package models

import (
	"gorm.io/gorm"
)

type Unit struct {
	Id          int     `gorm:"primaryKey"`
	Name        string  `gorm:"index;not null"`
	Symbol      string  `gorm:"index;not null"`
	Description *string `gorm:"default:null"`
	HasShopId
	HasIsActive
}

type UnitResource struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type UnitDetailResource struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Symbol      string  `json:"symbol"`
	Description *string `json:"description"`
	HasShopId
	HasIsActive
}

// clean related cache when category is updated
func (readers *ReadServices) AfterUnitUpdate(db *gorm.DB, shopId string, id int) error {

	if err := readers.UnitGetService.CleanCache(id); err != nil {
		return err
	}
	if err := readers.UnitListService.CleanCache(shopId); err != nil {
		return err
	}
	var relatedItemIds []int
	if err := db.Model(&Item{}).Where("unit_id = ?", id).Pluck("id", &relatedItemIds).Error; err != nil {
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
