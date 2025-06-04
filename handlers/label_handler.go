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
		return vr.CreateFormLabel(w, ds.UserId)
	})
}

func ShowLabelEdit(vr *views.Renderer, db *gorm.DB) http.HandlerFunc {
	return ResourceHandler(func(w http.ResponseWriter, r *http.Request, s *Session) error {
		var label models.Label
		if err := db.First(&label, s.ResId).Error; err != nil {
			return err
		}
		return vr.EditFormLabel(w, s.UserId, s.ResId, &label)
	})
}

func ShowLabelIndex(vr *views.Renderer, db *gorm.DB) http.HandlerFunc {
	return DefaultHandler(func(w http.ResponseWriter, r *http.Request, s *DefaultSession) error {
		var results []models.Label
		if err := db.Find(&results).Error; err != nil {
			return err
		}
		return vr.IndexPageLabel(w, results, r.Header.Get("HX-Request") == "")
	})
}

func HandleLabelUpdate(vr *views.Renderer, db *gorm.DB) http.HandlerFunc {
	return ResourceHandler(func(w http.ResponseWriter, r *http.Request, s *Session) error {
		input, errmap := ParseLabel(w, r)
		if len(errmap) > 0 {
			w.WriteHeader(http.StatusNoContent)
			return vr.EditFormLabelWithErrors(w, s.UserId, s.ResId, input, errmap)
		}

		// if m := input.Validate(container.db, container.resId, container.userId); len(m) > 0 {
		// 	return &input, m
		// }
		var label models.Label
		if err := db.First(&label, s.ResId).Error; err != nil {
			return err
		}
		updates := map[string]any{
			"Name":        input.Name,
			"Description": input.Description,
		}
		if err := db.Model(&label).Updates(updates).Error; err != nil {
			return err
		}
		return vr.SuccessLabelEdit(w, &label)
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

	return DefaultHandler(func(w http.ResponseWriter, r *http.Request, ds *DefaultSession) error {
		input, errmap := ParseLabel(w, r)
		if len(errmap) > 0 {
			w.WriteHeader(http.StatusNoContent)
			return vr.CreateFormLabelWithErrors(w, ds.UserId, input, errmap)
		}
		if err := db.Create(&input).Error; err != nil {
			return err
		}

		return vr.SuccessLabelCreate(w, input)
	})
}

func ParseLabel(w http.ResponseWriter, r *http.Request) (*models.Label, map[string]error) {

	var input models.Label

	scans := [...]ScanFormFunc{
		newScanner(r, "name", &input.Name, formscanner.StringRequired, formscanner.MinMax(4, 20)),
		newScanner(r, "description", &input.Description, formscanner.String, formscanner.MinMax(4, 20)),
	}

	m := runScanners(scans[:])

	return &input, m
}
