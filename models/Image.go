package models

type Image struct {
	Id            int `gorm:"primaryKey"`
	Url           string
	ReferenceId   int                `gorm:"index;not null"`
	ReferenceType ImageReferenceType `gorm:"index;not null;enum('ghost');default:ghost"`
	Size          int64              `gorm:"not null"`
}
