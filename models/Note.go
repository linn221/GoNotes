package models

import (
	"context"
	"errors"
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

	if input.Reminder.Before(time.Now()) {
		return errors.New("remind date should be in the future")
	}
	userFilter := NewFilter("user_id = ?", userId)
	return Validate(db, NewExistsRule("labels", input.LabelId, "label not found", userFilter),
		NewUniqueRule("notes", "title", input.Title, id, "duplicate title", userFilter),
	)
}

type NoteResource struct {
	Id           int
	Title        string
	Description  string
	Body         string
	LabelId      int
	LabelName    string
	ParentNoteId int
	Reminder     string
}

type NoteService struct {
	db *gorm.DB
}

func (s *NoteService) fetch(ctx context.Context, userId int, id int) (*Note, error) {
	return first[Note](s.db.WithContext(ctx), userId, id)
}

func (s *NoteService) Create(ctx context.Context, userId int, input *Note) (*Note, error) {

	if err := input.validate(s.db.WithContext(ctx), userId, 0); err != nil {
		return nil, err
	}
	err := s.db.WithContext(ctx).Create(&input).Error
	return input, err
}

func (s *NoteService) Update(ctx context.Context, userId int, id int, input *Note) (*Note, error) {

	if err := input.validate(s.db.WithContext(ctx), userId, id); err != nil {
		return nil, err
	}
	note, err := s.fetch(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	updates := map[string]any{
		"Title":       input.Title,
		"Description": input.Description,
		"Body":        input.Body,
		"LabelId":     input.LabelId,
		"Reminder":    input.Reminder,
	}
	if err := s.db.WithContext(ctx).Model(&note).Updates(updates).Error; err != nil {
		return nil, err
	}
	return note, nil
}

func (s *NoteService) UpdateBody(ctx context.Context, userId int, id int, body string) (*Note, error) {
	note, err := s.fetch(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	err = s.db.WithContext(ctx).Model(&note).UpdateColumn("body", body).Error
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (s *NoteService) Delete(ctx context.Context, userId int, id int) (*Note, error) {
	note, err := s.fetch(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Delete(&note).Error; err != nil {
		return nil, err
	}
	return note, nil
}

func (s *NoteService) convertToResoruce(note *Note) *NoteResource {
	res := NoteResource{
		Id:          note.Id,
		Title:       note.Title,
		Description: note.Description,
		Body:        note.Body,
		LabelId:     note.LabelId,
		LabelName:   note.Label.Name,
		Reminder:    note.Reminder.Format(time.DateOnly),
	}
	return &res
}

func (s *NoteService) Get(ctx context.Context, userId int, id int) (*NoteResource, error) {

	note, err := first[Note](s.db.WithContext(ctx), userId, id, "Label")
	if err != nil {
		return nil, err
	}
	res := s.convertToResoruce(note)
	return res, nil
}

func (s *NoteService) ListNotes(ctx context.Context, userId int) ([]*NoteResource, error) {
	notes, err := find[Note](s.db.WithContext(ctx), userId, "Label")
	if err != nil {
		return nil, err
	}
	resCollection := make([]*NoteResource, 0, len(notes))
	for _, n := range notes {
		resCollection = append(resCollection, s.convertToResoruce(&n))
	}

	return resCollection, nil
}
