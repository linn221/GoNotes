package models

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"linn221/shop/utils"
	"net/http"
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
	RemindDate   time.Time `gorm:"index;default:null"`
	HasUserId
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (input *Note) validate(db *gorm.DB, userId int, id int) error {

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
	RemindDate   MyDate
	CreatedAt    MyDateTime
	UpdatedAt    MyDateTime
}

type NoteService struct {
	db *gorm.DB
}

func (s *NoteService) fetch(ctx context.Context, userId int, id int, preloads ...string) (*Note, error) {
	return first[Note](s.db.WithContext(ctx), userId, id, preloads...)
}

func (s *NoteService) Create(ctx context.Context, userId int, input *Note) (*Note, error) {

	if err := input.validate(s.db.WithContext(ctx), userId, 0); err != nil {
		return nil, err
	}
	input.UserId = userId

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
		"RemindDate":  input.RemindDate,
	}
	if err := s.db.WithContext(ctx).Model(&note).Updates(updates).Error; err != nil {
		return nil, err
	}
	return note, nil
}

func (s *NoteService) UpdateBody(ctx context.Context, userId int, id int, body string) (*NoteResource, error) {
	note, err := s.fetch(ctx, userId, id, "Label")
	if err != nil {
		return nil, err
	}
	err = s.db.WithContext(ctx).Model(&note).UpdateColumn("body", body).Error
	if err != nil {
		return nil, err
	}
	return s.ConvertToResource(note), nil
}

func (s *NoteService) UpdateLabel(ctx context.Context, userId int, id int, labelId int) (*NoteResource, error) {
	note, err := s.fetch(ctx, userId, id)
	if err != nil {
		return nil, err
	}

	if err := Validate(s.db.WithContext(ctx), NewExistsRule("labels", labelId, "label not found", NewFilter("user_id = ? AND is_active = 1", userId))); err != nil {
		return nil, err
	}

	err = s.db.WithContext(ctx).Model(&note).UpdateColumn("label_id", labelId).Error
	if err != nil {
		return nil, err
	}

	note, err = s.fetch(ctx, userId, id, "Label")
	if err != nil {
		return nil, err
	}
	return s.ConvertToResource(note), nil
}

func (s *NoteService) UpdateRemindDate(ctx context.Context, userId int, id int, inputdate time.Time) (*NoteResource, error) {
	note, err := s.fetch(ctx, userId, id)
	if err != nil {
		return nil, err
	}

	if inputdate.Before(time.Now()) {
		return nil, errors.New("remind date must be in future")
	}

	err = s.db.WithContext(ctx).Model(&note).UpdateColumn("remind_date", inputdate).Error
	if err != nil {
		return nil, err
	}

	note, err = s.fetch(ctx, userId, id, "Label")
	if err != nil {
		return nil, err
	}
	return s.ConvertToResource(note), nil
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

func (s *NoteService) ConvertToResource(note *Note) *NoteResource {
	var remindDate MyDate
	remindDate = MyDate{note.RemindDate}
	res := NoteResource{
		Id:          note.Id,
		Title:       note.Title,
		Description: note.Description,
		Body:        note.Body,
		LabelId:     note.LabelId,
		LabelName:   note.Label.Name,
		RemindDate:  remindDate,
		CreatedAt:   MyDateTime{note.CreatedAt.Local()},
		UpdatedAt:   MyDateTime{note.UpdatedAt.Local()},
	}
	return &res
}

func (s *NoteService) Get(ctx context.Context, userId int, id int) (*NoteResource, error) {

	note, err := first[Note](s.db.WithContext(ctx), userId, id, "Label")
	if err != nil {
		return nil, err
	}
	res := s.ConvertToResource(note)
	return res, nil
}

// func (s *NoteService) GetRemindNotes(ctx context.Context, userId int) ([]Note, error) {
// }

type NoteSearchParam struct {
	LabelId int
}

func (s *NoteService) listAllNotes(ctx context.Context, userId int) ([]Note, error) {
	notes, err := find[Note](s.db.WithContext(ctx).Order("updated_at DESC"), userId, "Label")
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (s *NoteService) ListNotes(ctx context.Context, userId int, param *NoteSearchParam) ([]*NoteResource, error) {
	var notes []Note
	var err error
	if param != nil {
		dbCtx := s.db.WithContext(ctx).Preload("Label").Where("user_id = ?", userId)
		if param.LabelId > 0 {
			dbCtx.Where("label_id = ?", param.LabelId)
		}

		err = dbCtx.Order("updated_at DESC").Find(&notes).Error
	} else {
		notes, err = s.listAllNotes(ctx, userId)
	}

	if err != nil {
		return nil, err
	}
	resCollection := make([]*NoteResource, 0, len(notes))
	for _, n := range notes {
		resCollection = append(resCollection, s.ConvertToResource(&n))
	}

	return resCollection, nil
}

func (n *NoteResource) Checksum() string {
	return utils.HashString(fmt.Sprintf("title=%s&remindDate=%s&label=%s&desc=%s&body=%s",
		n.Title,
		n.RemindDate.Format(time.DateOnly),
		n.LabelName,
		n.Description,
		n.Body,
	))
}

func (n *NoteResource) CsvValues() []string {
	values := make([]string, 0, 6)
	values = append(values,
		n.Title,
		n.RemindDate.Format(time.DateOnly),
		n.LabelName,
		n.Description,
		n.Body,
		n.Checksum(),
	)
	return values
}

func (s *NoteService) Export(ctx context.Context, w http.ResponseWriter, notes []*NoteResource) error {

	filename := time.Now().Format(time.DateOnly) + ".csv"
	// Set headers to force download
	w.Header().Set("Content-Disposition", "attachment;filename="+filename)
	w.Header().Set("Content-Type", "text/csv")
	csvWriter := csv.NewWriter(w)
	if err := csvWriter.Write([]string{"title", "remind_date", "label", "description", "body", "checksum"}); err != nil {
		return err
	}
	for _, note := range notes {
		csvWriter.Write(note.CsvValues())
	}
	csvWriter.Flush()
	return csvWriter.Error()
}

type ImportedNote struct {
	Title       string
	RemindDate  time.Time
	Label       string
	Description string
	Body        string
	Checksum    string
}

func parseImportedNote(row []string) *ImportedNote {
	t, err := time.Parse(time.DateOnly, row[1])
	if err != nil {
		panic(err)
	}
	return &ImportedNote{
		Title:       row[0],
		RemindDate:  t,
		Label:       row[2],
		Description: row[3],
		Body:        row[4],
		Checksum:    row[5],
	}
}

func (s *NoteService) ImportNotes(ctx context.Context, userId int, rows [][]string) error {

	labelService := LabelService{db: s.db}
	labels, err := labelService.ListAll(ctx, userId)
	if err != nil {
		return err
	}

	labelMap := make(map[string]int, len(labels))
	for _, label := range labels {
		labelMap[label.Name] = label.Id
	}
	existingNotes, err := s.ListNotes(ctx, userId, nil)
	if err != nil {
		return err
	}
	existingChecksums := make(map[string]struct{}, len(existingNotes))
	for _, note := range existingNotes {
		existingChecksums[note.Checksum()] = struct{}{}
	}

	tx := s.db.WithContext(ctx).Begin()
	labelService.db = tx
	for _, row := range rows[1:] {
		imported := parseImportedNote(row)
		labelId, labelFound := labelMap[imported.Label]
		if !labelFound {
			// create labe if not found
			newlabel, err := labelService.Create(ctx, userId, &Label{Name: imported.Label})
			if err != nil {
				return err
			}
			labelId = newlabel.Id
			labelMap[newlabel.Name] = newlabel.Id
		}

		// if the note exits or not
		_, exists := existingChecksums[imported.Checksum]
		if !exists {
			newNote := Note{
				Title:       imported.Title,
				Description: imported.Description,
				LabelId:     labelId,
				RemindDate:  imported.RemindDate,
				Body:        imported.Body,
			}
			newNote.UserId = userId
			if err := tx.Create(&newNote).Error; err != nil {
				return err
			}
			existingChecksums[imported.Checksum] = struct{}{}
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
