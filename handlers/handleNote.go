package handlers

import (
	"context"
	"encoding/csv"
	"io"
	"linn221/shop/formscanner"
	"linn221/shop/models"
	"linn221/shop/services"
	"linn221/shop/views"
	"net/http"
	"strconv"
)

func parseNote(r *http.Request) (*models.Note, services.FormErrors) {
	var input models.Note
	scans := [...]ScannerFunc{
		newScannerFunc(r, "title", &input.Title, formscanner.StringRequired, formscanner.MinMax(4, 255)),
		newScannerFunc(r, "description", &input.Description, formscanner.String, formscanner.Max(500)),
		newScannerFunc(r, "body", &input.Body, formscanner.String),
		newScannerFunc(r, "label_id", &input.LabelId, formscanner.IntRequired, formscanner.Gte(1)),
		newScannerFunc(r, "remind", &input.RemindDate, formscanner.Date, formscanner.InFuture),
	}

	fe := runScanners(scans[:])

	return &input, fe
}

func ShowNoteCreate(t *views.Templates, labelService *models.LabelService) http.HandlerFunc {
	return DefaultHandler(t, func(ctx context.Context, r *http.Request, session *DefaultSession, vr *views.Renderer) error {
		labels, err := labelService.ListActiveOnly(ctx, session.UserId)
		if err != nil {
			return err
		}
		labelId, ok := getQueryInt(r, "label_id")
		if !ok {
			labelId = labels[0].Id
		}
		return vr.NoteCreateForm(session.UserId, labels, labelId)
	})
}

func getQueryInt(r *http.Request, key string) (int, bool) {
	s := r.URL.Query().Get(key)
	if s != "" {
		i, _ := strconv.Atoi(s)
		if i > 0 {
			return i, true
		}
	}
	return 0, false
}

func ShowNoteEdit(t *views.Templates, noteService *models.NoteService, labelService *models.LabelService) http.HandlerFunc {
	return ResourceHandler(t, func(ctx context.Context, r *http.Request, session *Session, vr *views.Renderer) error {
		res, err := noteService.Get(ctx, session.UserId, session.ResId)
		if err != nil {
			return err
		}
		labels, err := labelService.ListActiveOnly(ctx, session.UserId)
		if err != nil {
			return err
		}

		return vr.NoteEditForm(session.UserId, session.ResId, res, labels)
	})
}

func ShowNoteIndex(t *views.Templates, noteService *models.NoteService) http.HandlerFunc {
	parseSearchParam := func(r *http.Request) *models.NoteSearchParam {
		var searchParam *models.NoteSearchParam
		labelId, ok := getQueryInt(r, "label_id")
		if ok {
			searchParam = &models.NoteSearchParam{LabelId: labelId}
		}
		return searchParam
	}
	return DefaultHandler(t, func(ctx context.Context, r *http.Request, session *DefaultSession, vr *views.Renderer) error {
		//parse search param
		searchParam := parseSearchParam(r)
		notes, err := noteService.ListNotes(ctx, session.UserId, searchParam)
		if err != nil {
			return err
		}
		return vr.NoteIndexPage(notes)
	})
}

func HandleNoteCreate(t *views.Templates, noteService *models.NoteService, labelService *models.LabelService) http.HandlerFunc {

	return CreateHandler(t, parseNote, func(w http.ResponseWriter, r *http.Request, session *DefaultSession, input *models.Note, fe services.FormErrors, vr *views.Renderer) error {
		if len(fe) > 0 {
			labels, err := labelService.ListActiveOnly(r.Context(), session.UserId)
			if err != nil {
				return err
			}
			return vr.NoteCreateError(input, labels, fe)
		}
		if input.Description == "" {
			input.Description = input.Title
		}
		_, err := noteService.Create(r.Context(), session.UserId, input)
		if err != nil {
			return err
		}

		htmxRedirect(w, "/notes")
		return nil
	})
}

func HandleNoteUpdate(t *views.Templates, noteService *models.NoteService, labelService *models.LabelService) http.HandlerFunc {
	return UpdateHandler(t, parseNote, func(w http.ResponseWriter, r *http.Request, session *Session, input *models.Note, fe services.FormErrors, renderer *views.Renderer) error {
		if len(fe) > 0 {
			labels, err := labelService.ListActiveOnly(r.Context(), session.UserId)
			if err != nil {
				return err
			}

			return renderer.NoteEditError(session.UserId, session.ResId, input, fe, labels)
		}
		_, err := noteService.Update(r.Context(), session.UserId, session.ResId, input)
		if err != nil {
			return err
		}
		htmxRedirect(w, "/notes")
		return nil
	})
}

func HandleNoteDelete(noteService *models.NoteService) http.HandlerFunc {
	return DeleteHandler(func(ctx context.Context, r *http.Request, userId, resId int) error {
		_, err := noteService.Delete(ctx, userId, resId)
		if err != nil {
			return err
		}

		return nil
	})
}

func HandleNoteUpdateBody(t *views.Templates, noteService *models.NoteService) http.HandlerFunc {
	return ResourceHandler(t, func(ctx context.Context, r *http.Request, session *Session, vr *views.Renderer) error {
		body := r.PostFormValue("body")
		updated, err := noteService.UpdateBody(r.Context(), session.UserId, session.ResId, body)
		if err != nil {
			return err
		}
		return vr.NoteUpdateBodySuccess(updated)
	})
}

func HandleNoteExport(noteService *models.NoteService) http.HandlerFunc {
	return MinHandler(func(w http.ResponseWriter, r *http.Request, userId int) error {
		ctx := r.Context()
		notes, err := noteService.ListNotes(ctx, userId, nil)
		if err != nil {
			return err
		}
		return noteService.Export(ctx, w, notes)
	})
}

func HandleNoteImport(noteService *models.NoteService) http.HandlerFunc {

	return MinHandler(func(w http.ResponseWriter, r *http.Request, userId int) error {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			return err
		}

		// Get the file from form field "csvfile"
		file, _, err := r.FormFile("csvfile")
		if err != nil {
			return err
		}
		defer file.Close()

		// Parse the CSV file
		reader := csv.NewReader(file)
		var records [][]string

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			records = append(records, record)
		}
		return noteService.ImportNotes(r.Context(), userId, records)
	})
}
