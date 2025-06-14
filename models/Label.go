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

func (s *LabelService) fetch(ctx context.Context, userId int, id int) (*Label, error) {
	return first[Label](s.db.WithContext(ctx), userId, id)
}

func (input *Label) validate(db *gorm.DB, userId int, id int) error {
	if err := Validate(db, NewUniqueRule("labels", "name", input.Name, id, "duplicate label name", NewFilter("user_id = ?", userId))); err != nil {
		return err
	}
	return nil
}

func (s *LabelService) Create(ctx context.Context, userId int, input *Label) (*Label, error) {
	if err := input.validate(s.db.WithContext(ctx), userId, 0); err != nil {
		return nil, err
	}

	input.UserId = userId
	err := s.db.WithContext(ctx).Create(&input).Error
	if err != nil {
		return nil, err
	}
	return input, nil
}

func (s *LabelService) Update(ctx context.Context, userId int, id int, input *Label) (*Label, error) {
	if err := input.validate(s.db.WithContext(ctx), userId, id); err != nil {
		return nil, err
	}
	updates := map[string]any{
		"Name":        input.Name,
		"Description": input.Description,
	}
	label, err := s.fetch(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	err = s.db.WithContext(ctx).Model(&label).Updates(updates).Error

	return label, err
}

func (s *LabelService) Delete(ctx context.Context, userId int, id int) (*Label, error) {
	if err := Validate(s.db.WithContext(ctx),
		NewNoResultRule("notes", "cannot delete label that has been used in note. toggle inactive if you don't want to see",
			NewFilter("label_id = ?", id))); err != nil {
		return nil, err
	}
	label, err := s.fetch(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	err = s.db.WithContext(ctx).Delete(&label).Error
	return label, err
}

func (s *LabelService) ToggleActive(ctx context.Context, userId int, id int, isActive bool) (*Label, error) {
	label, err := s.fetch(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	if label.IsActive != isActive {
		err = s.db.WithContext(ctx).Model(&label).UpdateColumn("is_active", isActive).Error
		if err != nil {
			return nil, err
		}
	}
	return label, nil
}

func (s *LabelService) Get(ctx context.Context, userId int, id int) (*Label, error) {

	return s.fetch(ctx, userId, id)
}

func (s *LabelService) ListAll(ctx context.Context, userId int) ([]Label, error) {

	var labels []Label
	err := s.db.WithContext(ctx).Where("user_id = ?", userId).Find(&labels).Error
	return labels, err
}

func (s *LabelService) ListActiveOnly(ctx context.Context, userId int) ([]Label, error) {

	var labels []Label
	err := s.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userId, true).Find(&labels).Error
	return labels, err
}
