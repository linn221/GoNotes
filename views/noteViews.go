package views

import (
	"linn221/shop/models"
	"time"
)

func (r *Renderer) ShowNoteCreate(userId int, labels []models.Label, labelId int) error {
	return r.templates.noteCreateTemplate.Execute(r.w, map[string]any{
		"Labels":  labels,
		"LabelId": labelId,
	})
}

func (r *Renderer) HandleNotePartialUpdate(note *models.NoteResource, timezone string) error {
	type Note struct {
		*models.NoteResource
		ExpandNote bool
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	note.CreatedAt.Time = note.CreatedAt.In(loc)
	note.UpdatedAt.Time = note.UpdatedAt.In(loc)
	note.RemindDate.Time = note.RemindDate.In(loc)

	return r.templates.noteTemplate.ExecuteTemplate(r.w, "note", Note{NoteResource: note, ExpandNote: true})
}

func (r *Renderer) ShowNoteEdit(userId int, resId int, res *models.NoteResource, labels []models.Label) error {
	return r.templates.noteEditTemplate.Execute(r.w, map[string]any{
		"Id":     resId,
		"Res":    res,
		"Labels": labels,
	})
}

func (r *Renderer) ShowNotePartialEditLabel(res *models.NoteResource, labels []models.Label) error {
	return r.templates.noteTemplate.ExecuteTemplate(r.w, "edit-label", map[string]any{
		"Labels": labels,
		"Res":    res,
	})
}

func (r *Renderer) ShowNotePartialEditBody(res *models.NoteResource) error {
	return r.templates.noteTemplate.ExecuteTemplate(r.w, "edit-body", res)
}

func (r *Renderer) RenderNoteIndex(notes []*models.NoteResource, labels []models.Label, timezone string) error {
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

func (r *Renderer) ShowNoteImport() error {
	return r.templates.importNoteTemplate.Execute(r.w, nil)
}
