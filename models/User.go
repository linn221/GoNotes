package models

type User struct {
	Id       int    `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique"`
	PhoneNo  string `gorm:"unique"`
	Password string `gorm:"index;not null"`
	HasIsActive
	HasShopId
}
