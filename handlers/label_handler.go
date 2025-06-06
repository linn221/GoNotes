package handlers

import (
	"linn221/shop/formscanner"
	"linn221/shop/models"
	"linn221/shop/views"
	"net/http"

	"gorm.io/gorm"
)

func ShowLabelCreate(vr *views.Renderer) http.HandlerFunc {
	return DefaultHandler(func(w http.ResponseWriter, r *http.Request, ds *DefaultSession) error {
		return vr.LabelCreateForm(w, ds.UserId)
	})
}

func ShowLabelEdit(vr *views.Renderer, db *gorm.DB) http.HandlerFunc {
	return ResourceHandler(func(w http.ResponseWriter, r *http.Request, s *Session) error {
		var label models.Label
		if err := db.First(&label, s.ResId).Error; err != nil {
			return err
		}
		return vr.LabelEditForm(w, s.UserId, s.ResId, &label)
	})
}

func ShowLabelIndex(vr *views.Renderer, db *gorm.DB) http.HandlerFunc {
	return DefaultHandler(func(w http.ResponseWriter, r *http.Request, s *DefaultSession) error {
		var results []models.Label
		if err := db.Find(&results).Error; err != nil {
			return err
		}
		return vr.LabelIndexPage(w, results, r.Header.Get("HX-Request") == "")
	})
}

func HandleLabelUpdate(vr *views.Renderer, db *gorm.DB) http.HandlerFunc {
	return UpdateHandler(UpdateResource[models.Label]{
		parseInput: parseLabel,
		handleParseError: func(w http.ResponseWriter, r *http.Request, s *Session, l *models.Label, fe formErrors) error {
			return vr.LabelUpdateError(w, s.UserId, s.ResId, l, fe)
		},
		handle: func(w http.ResponseWriter, r *http.Request, s *Session, l *models.Label) error {
			var label models.Label
			ctx := r.Context()
			if err := db.WithContext(ctx).First(&label, s.ResId).Error; err != nil {
				return err
			}
			updates := map[string]any{
				"Name":        l.Name,
				"Description": l.Description,
			}
			if err := db.WithContext(ctx).Model(&label).Updates(updates).Error; err != nil {
				return err
			}
			return vr.LabelUpdateOk(w, &label)
		},
	})
}

func HandleLabelDelete(db *gorm.DB) http.HandlerFunc {
	return ResourceHandler(func(w http.ResponseWriter, r *http.Request, s *Session) error {
		var label models.Label
		if err := db.First(&label, s.ResId).Error; err != nil {
			return err
		}
		if err := db.Delete(&label).Error; err != nil {
			return err
		}
		w.WriteHeader(http.StatusOK)
		return nil
	})
}

func HandleLabelCreate(vr *views.Renderer, db *gorm.DB) http.HandlerFunc {

	return CreateHandler(CreateResource[models.Label]{
		parseInput: parseLabel,
		handleParseError: func(w http.ResponseWriter, r *http.Request, ds *DefaultSession, l *models.Label, fe formErrors) error {
			return vr.LabelCreateError(w, ds.UserId, l, fe)
		},
		handle: func(w http.ResponseWriter, r *http.Request, ds *DefaultSession, l *models.Label) error {

			if err := db.WithContext(r.Context()).Create(&l).Error; err != nil {
				return err
			}
			return vr.LabelCreateOk(w, l)
		},
	})
}

func parseLabel(r *http.Request) (*models.Label, formErrors) {

	var input models.Label

	scans := [...]ScannerFunc{
		newScannerFunc(r, "name", &input.Name, formscanner.StringRequired, formscanner.MinMax(4, 20)),
		newScannerFunc(r, "description", &input.Description, formscanner.String, formscanner.MinMax(4, 20)),
	}

	m := runScanners(scans[:])

	return &input, m
}
