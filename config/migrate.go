package config

import (
	"linn221/shop/models"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {

	err := db.AutoMigrate(&models.Shop{}, &models.User{}, &models.Image{},
		&models.Label{})
	if err != nil {
		panic("Error migrating: " + err.Error())
	}
}
