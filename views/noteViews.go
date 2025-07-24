package views

import (
	"fmt"
	"linn221/shop/formscanner"
	"linn221/shop/models"
	"linn221/shop/services"
	"linn221/shop/utils"
)

func (r *Renderer) NoteCreateForm(userId int, labels []models.Label, labelId int) error {
	return r.templates.noteCreateTemplate.Execute(r.w, map[string]any{
		"Labels":  labels,
		"LabelId": labelId,
	})
}

func (r *Renderer) NoteCreateError(input *models.Note, labels []models.Label, errmap services.FormErrors) error {
	m := map[string]any{
		"Labels": labels,

		"title":       NewFormInput(input.Title, errmap["title"]),
		"description": NewFormInput(input.Description, errmap["description"]),
		"body":        NewFormInput(input.Body, errmap["body"]),
		"label_id":    NewFormInput(fmt.Sprint(input.LabelId), errmap["label_id"]),
		"remind":      NewFormInput(input.RemindDate.Format(formscanner.MyDateFormat), errmap["remind"]),
	}
	return r.templates.noteCreateTemplate.ExecuteTemplate(r.w, "create_form2", m)
}

func (r *Renderer) NoteUpdateBodySuccess(note *models.NoteResource, labels []models.Label, timezone string) error {

	noteData := newNoteData(note)
	return r.templates.noteTemplate.ExecuteTemplate(r.w, "note", map[string]any{
		"Res":      noteData,
		"Labels":   labels,
		"Timezone": timezone,
	})
}

func (r *Renderer) NoteEditForm(userId int, resId int, res *models.NoteResource, labels []models.Label) error {
	return r.templates.noteEditTemplate.Execute(r.w, map[string]any{
		"Id":     resId,
		"Res":    res,
		"Labels": labels,
	})
}

func (r *Renderer) NoteEditError(userId int, resId int, input *models.Note, errmap services.FormErrors, labels []models.Label) error {
	m := map[string]any{
		"Labels": labels,
		"Id":     resId,

		"title":       NewFormInput(input.Title, errmap["title"]),
		"description": NewFormInput(input.Description, errmap["description"]),
		"body":        NewFormInput(input.Body, errmap["body"]),
		"label_id":    NewFormInput(fmt.Sprint(input.LabelId), errmap["label_id"]),
		"remind":      NewFormInput(input.RemindDate.Format(formscanner.MyDateFormat), errmap["remind"]),
	}
	return r.templates.noteEditTemplate.ExecuteTemplate(r.w, "edit_form2", m)
}

//	func (r Renderer) NoteCreateSuccess(note *models.NoteResource) error {
//		return r.templates.noteTemplate.ExecuteTemplate(r.w, "note", note)
//	}
type noteData struct {
	*models.NoteResource
	BodyShort        string
	ReadMoreRequired bool
	HasBody          bool
	Labels           []models.Label
}

func newNoteData(note *models.NoteResource) *noteData {
	excerpt, readMoreRequired := utils.GenerateExcerpt(note.Body, 20)
	return &noteData{
		NoteResource:     note,
		BodyShort:        excerpt,
		HasBody:          note.Body != "",
		ReadMoreRequired: readMoreRequired,
	}
}

func (r *Renderer) NoteIndexPage(notes []*models.NoteResource, labels []models.Label, timezone string) error {
	noteCollection := make([]*noteData, 0, len(notes))
	for _, note := range notes {
		noteCollection = append(noteCollection, newNoteData(note))
	}

	return r.templates.noteTemplate.Execute(r.w, map[string]any{
		"ResList":   noteCollection,
		"Labels":    labels,
		"PageTitle": "Notes",
		"Nav":       NavNotes,
		"Timezone":  timezone,
	})
}

func (r *Renderer) NoteImportPage() error {
	return r.templates.importNoteTemplate.Execute(r.w, nil)
}
