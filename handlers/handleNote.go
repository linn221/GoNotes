package handlers

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"linn221/shop/formscanner"
	"linn221/shop/models"
	"linn221/shop/services"
	"linn221/shop/views"
	"net/http"
	"strconv"
	"time"
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
func ShowNoteImport(t *views.Templates) http.HandlerFunc {
	return DefaultHandler(t, func(ctx context.Context, r *http.Request, session *DefaultSession, vr *views.Renderer) error {
		return vr.ShowNoteImport()
	})
}
func ShowNoteCreate(t *views.Templates, labelService *models.LabelService) http.HandlerFunc {
	return DefaultHandler(t, func(ctx context.Context, r *http.Request, session *DefaultSession, vr *views.Renderer) error {
		labels, err := labelService.ListActiveOnly(ctx, session.UserId)
		if err != nil {
			return err
		}
		if len(labels) <= 0 {
			return errors.New("please create a label to continue")
		}
		labelId, ok := getQueryInt(r, "label_id")
		if !ok {
			labelId = labels[0].Id
		}
		return vr.ShowNoteCreate(session.UserId, labels, labelId)
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

		return vr.ShowNoteEdit(session.UserId, session.ResId, res, labels)
	})
}

func ShowNotePartialEdit(t *views.Templates, noteService *models.NoteService, labelService *models.LabelService, tz func(context.Context) string) http.HandlerFunc {
	return ResourceHandler(t, func(ctx context.Context, r *http.Request, session *Session, vr *views.Renderer) error {
		res, err := noteService.Get(ctx, session.UserId, session.ResId)
		if err != nil {
			return err
		}
		part := r.URL.Query().Get("part")
		if part == "" {
			return errors.New("part must be specified")
		}
		switch part {
		case "body":
			return vr.ShowNotePartialEditBody(res)
		case "label":
			labels, err := labelService.ListActiveOnly(ctx, session.UserId)
			if err != nil {
				return err
			}
			return vr.ShowNotePartialEditLabel(res, labels)
		case "none":
			return vr.HandleNotePartialUpdate(res, tz(ctx))
		default:
			return errors.New("invalid query")
		}
	})
}

func RenderNoteIndex(t *views.Templates, noteService *models.NoteService, labelService *models.LabelService, tz func(ctx context.Context) string) http.HandlerFunc {
	parseSearchParam := func(r *http.Request) models.NoteSearchParam {
		var searchParam models.NoteSearchParam
		searchParam.LabelId, _ = getQueryInt(r, "label_id")
		searchParam.Search = r.URL.Query().Get("search")
		return searchParam
	}
	return DefaultHandler(t, func(ctx context.Context, r *http.Request, session *DefaultSession, vr *views.Renderer) error {
		//parse search param
		timezone := tz(ctx)
		searchParam := parseSearchParam(r)
		notes, err := noteService.ListNotes(ctx, session.UserId, searchParam)
		if err != nil {
			return err
		}
		labels, err := labelService.ListActiveOnly(ctx, session.UserId)
		if err != nil {
			return err
		}
		return vr.RenderNoteIndex(notes, labels, timezone)
	})
}

func HandleNoteCreate(t *views.Templates, noteService *models.NoteService, labelService *models.LabelService) http.HandlerFunc {

	return CreateHandler(t, parseNote, func(w http.ResponseWriter, r *http.Request, session *DefaultSession, input *models.Note, vr *views.Renderer) error {
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
	return UpdateHandler(t, parseNote, func(w http.ResponseWriter, r *http.Request, session *Session, input *models.Note, renderer *views.Renderer) error {
		_, err := noteService.Update(r.Context(), session.UserId, session.ResId, input)
		if err != nil {
			return err
		}
		htmxRedirect(w, "/notes")
		return nil
	})
}

func HandleNoteDelete(t *views.Templates, noteService *models.NoteService) http.HandlerFunc {
	return DeleteHandler(t, func(ctx context.Context, r *http.Request, userId, resId int) error {
		_, err := noteService.Delete(ctx, userId, resId)
		if err != nil {
			return err
		}

		return nil
	})
}

func HandleNotePartialUpdate(t *views.Templates, noteService *models.NoteService, labelService *models.LabelService, tz func(context.Context) string) http.HandlerFunc {
	return ResourceHandler(t, func(ctx context.Context, r *http.Request, session *Session, vr *views.Renderer) error {
		var updated *models.NoteResource
		var err error
		if body := r.PostFormValue("body"); body != "" {
			updated, err = noteService.UpdateBody(r.Context(), session.UserId, session.ResId, body)
		} else if labelIdStr := r.PostFormValue("label_id"); labelIdStr != "" {
			labelId, err2 := strconv.Atoi(labelIdStr)
			if err2 != nil {
				return err2
			}
			updated, err = noteService.UpdateLabel(ctx, session.UserId, session.ResId, labelId)
		} else if remindDateStr := r.PostFormValue("remind"); remindDateStr != "" {
			inputRemindDate, err2 := time.Parse(time.DateOnly, remindDateStr) // to avoid err being shadowed
			if err2 != nil {
				return err
			}
			updated, err = noteService.UpdateRemindDate(ctx, session.UserId, session.ResId, inputRemindDate)
		} else {
			err = errors.New("no form data")
		}

		if err != nil {
			return err
		}

		timezone := tz(ctx)
		return vr.HandleNotePartialUpdate(updated, timezone)
	})
}

func HandleNoteExport(noteService *models.NoteService) http.HandlerFunc {
	return MinHandler(func(w http.ResponseWriter, r *http.Request, userId int) error {
		ctx := r.Context()
		notes, err := noteService.ListNotes(ctx, userId, models.NoteSearchParam{})
		if err != nil {
			return err
		}
		err = noteService.Export(ctx, w, notes)
		if err != nil {
			return err
		}

		// htmxRedirect(w, "/")
		return nil
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
		err = noteService.ImportNotes(r.Context(), userId, records)
		if err != nil {
			return err
		}
		htmxRedirect(w, "/")
		return nil
	})
}
