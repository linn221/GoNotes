package models

import "time"

type Note struct {
	HasID
	Title        string `gorm:"index;not null"`
	Description  string `gorm:"index"`
	Body         string `gorm:"longtext"`
	LabelId      int    `gorm:"index;not null"`
	ParentNoteId int    `gorm:"index"`
	Label        Label
	Reminder     time.Time `gorm:"index"`
	HasUserId
}
