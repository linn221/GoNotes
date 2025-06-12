package models

import (
	"context"

	"gorm.io/gorm"
)

type Label struct {
	Id          int    `gorm:"primaryKey"`
	Name        string `gorm:"index;not null"`
	Description string
	HasIsActive
	HasUserId
}

type LabelService struct {
	db *gorm.DB
}

func (input *Label) validate(db *gorm.DB, userId int, id int) error {
	panic("2d")
}

func (s *LabelService) Create(ctx context.Context, userId int, input *Label) (*Label, error) {

	panic("2d")
}

func (s *LabelService) Update(ctx context.Context, userId int, id int, input *Label) (*Label, error) {

	panic("2d")
}

func (s *LabelService) Delete(ctx context.Context, userId int, id int) (*Label, error) {
	panic("2d")
}

func (s *LabelService) Get(ctx context.Context, userId int, id int) (*Label, error) {

	panic("2d")
}

func (s *LabelService) List(ctx context.Context, userId int) ([]Label, error) {

	panic("2d")
}
