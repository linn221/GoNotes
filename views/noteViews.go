package views

import (
	"linn221/shop/models"
	"linn221/shop/utils"
	"time"
)

type NoteData struct {
	*models.NoteResource
	ExpandNote bool
}

func (r *Renderer) ShowNoteCreate(userId int, labels []models.Label, labelId int) error {
	return r.templates.noteCreateTemplate.Execute(r.w,
		Page{
			PageTitle: "Create Note",
			Nav:       NavNotes,
			Data: H{
				"Labels":  labels,
				"LabelId": labelId,
			},
		})
}
func (r *Renderer) ShowNoteEdit(userId int, resId int, res *models.NoteResource, labels []models.Label) error {
	return r.templates.noteEditTemplate.Execute(r.w, Page{
		PageTitle: "Edit Note",
		Nav:       NavNotes,
		Res:       res,
		Data:      H{"LabelList": labels},
	})
}

func (r *Renderer) HandleNotePartialUpdate(note *models.NoteResource, timezone string) error {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}
	note.CreatedAt.Time = note.CreatedAt.In(loc)
	note.UpdatedAt.Time = note.UpdatedAt.In(loc)
	note.RemindDate.Time = note.RemindDate.In(loc)
	localDate, err := utils.GetLocalDate(timezone)
	if err != nil {
		return err
	}
	if note.RemindDate.Equal(localDate) {
		note.Pinned = true
	}

	return r.templates.noteTemplate.ExecuteTemplate(r.w, "note", ResourceData{Res: note, Data: H{"ExpandNote": true}})
}

func (r *Renderer) ShowNotePartialEditLabel(res *models.NoteResource, labels []models.Label) error {
	return r.templates.noteTemplate.ExecuteTemplate(r.w, "edit-label", ResourceData{Res: res, Data: H{"LabelList": labels}})
}

func (r *Renderer) ShowNotePartialEditBody(res *models.NoteResource) error {
	return r.templates.noteTemplate.ExecuteTemplate(r.w, "edit-body", ResourceData{Res: res})
}

func (r *Renderer) ShowNotePartialEditRemind(res *models.NoteResource) error {
	return r.templates.noteTemplate.ExecuteTemplate(r.w, "edit-remind", ResourceData{Res: res})

}
func (r *Renderer) RenderNoteIndex(notes []*models.NoteResource, labels []models.Label, timezone string) error {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	ResList := make([]NoteData, 0, len(notes))
	for _, note := range notes {
		note.CreatedAt.Time = note.CreatedAt.In(loc)
		note.UpdatedAt.Time = note.UpdatedAt.In(loc)
		note.RemindDate.Time = note.RemindDate.In(loc)
		ResList = append(ResList, NoteData{NoteResource: note})
	}

	return r.templates.noteTemplate.Execute(r.w, Page{
		PageTitle: "Notes",
		Nav:       NavNotes,
		Timezone:  timezone,
		ResList:   ResList,
		Data: H{
			"LabelList": labels,
		},
	})
}

func (r *Renderer) ShowNoteImport() error {
	return r.templates.importNoteTemplate.Execute(r.w, nil)
}
