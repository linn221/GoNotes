package models

import (
	"context"
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

func (input *Note) validate(db *gorm.DB, userId int, id int) error {

	panic("2d")
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

func (s *NoteService) Create(ctx context.Context, userId int, input *Note) (*Note, error) {

	panic("2d")
}

func (s *NoteService) Update(ctx context.Context, userId int, id int, input *Note) (*Note, error) {

	panic("2d")
}

func (s *NoteService) UpdateBody(ctx context.Context, userId int, id int, body string) (*Note, error) {
	panic("2d")
}

func (s *NoteService) Delete(ctx context.Context, userId int, id int) (*Note, error) {
	panic("2d")
}

func (s *NoteService) Get(ctx context.Context, userId int, id int) (*NoteResource, error) {

	panic("2d")
}

func (s *NoteService) ListNotes(ctx context.Context, userId int) ([]NoteResource, error) {

	panic("2d")
}
