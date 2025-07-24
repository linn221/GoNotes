package views

import (
	"linn221/shop/models"
	"time"
)

func (r *Renderer) NoteCreateForm(userId int, labels []models.Label, labelId int) error {
	return r.templates.noteCreateTemplate.Execute(r.w, map[string]any{
		"Labels":  labels,
		"LabelId": labelId,
	})
}

func (r *Renderer) NoteUpdateBodySuccess(note *models.NoteResource, labels []models.Label, timezone string) error {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	note.CreatedAt.Time = note.CreatedAt.In(loc)
	note.UpdatedAt.Time = note.UpdatedAt.In(loc)
	note.RemindDate.Time = note.RemindDate.In(loc)

	return r.templates.noteTemplate.ExecuteTemplate(r.w, "note", map[string]any{
		"Res":      note,
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

func (r *Renderer) NoteIndexPage(notes []*models.NoteResource, labels []models.Label, timezone string) error {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	for _, note := range notes {
		note.CreatedAt.Time = note.CreatedAt.In(loc)
		note.UpdatedAt.Time = note.UpdatedAt.In(loc)
		note.RemindDate.Time = note.RemindDate.In(loc)
	}

	return r.templates.noteTemplate.Execute(r.w, map[string]any{
		"ResList":   notes,
		"Labels":    labels,
		"PageTitle": "Notes",
		"Nav":       NavNotes,
		"Timezone":  timezone,
	})
}

func (r *Renderer) NoteImportPage() error {
	return r.templates.importNoteTemplate.Execute(r.w, nil)
}
