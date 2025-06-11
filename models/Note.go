package models

import (
	"time"

	"gorm.io/gorm"
)

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

type NoteResource struct {
	Id           int
	Title        string
	Description  string
	Body         string
	LabelId      int
	LabelName    string
	ParentNoteId int
	Reminder     time.Time
}

type NoteService struct {
	db *gorm.DB
}

func (service *NoteService) ListNotes(userId int) ([]NoteResource, error) {

	panic("2d")
}
