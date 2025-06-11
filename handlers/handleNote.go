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
