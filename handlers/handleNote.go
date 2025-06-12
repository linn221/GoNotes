package handlers

import (
	"linn221/shop/models"
	"linn221/shop/views"
	"net/http"
)

func HandleNoteIndex(t *views.Templates, getNotes func(int) ([]models.NoteResource, error)) http.HandlerFunc {
	// f := func(ctx context.Context, r *http.Request, session *DefaultSession, vr *views.Renderer) error {
	// 	notes, err := getNotes(session.UserId)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return nil
	// }
	return DefaultHandler(t, nil)
}

func HandleNoteCreate(t *views.Templates, noteService *models.NoteService) http.HandlerFunc {
	// return CreateHandler(t, nil, func(w http.ResponseWriter, r *http.Request, session *DefaultSession, input *models.Note, fe formErrors, vr *views.Renderer) error {
	// if len(fe) > 0 {

	// 	return nil
	// }

	// note, err := noteService.Create(r.Context(), session.UserId, input)
	// if err != nil {
	// 	return err
	// }

	// })
	return CreateHandler[models.Note](t, nil, nil)
}
