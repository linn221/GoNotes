package models

type Shop struct {
	Id      string `gorm:"primaryKey"`
	Name    string `gorm:"index;not null"`
	LogoUrl string `gorm:"index;not null"`
	Email   string `gorm:"index;not null"`
	PhoneNo string `gorm:"index;not null"`
	HasIsActive
}
