package handlers

import (
	"context"
	"linn221/shop/formscanner"
	"linn221/shop/models"
	"linn221/shop/views"
	"net/http"

	"gorm.io/gorm"
)

func parseLabel(r *http.Request) (*models.Label, formErrors) {

	var input models.Label

	scans := [...]ScannerFunc{
		newScannerFunc(r, "name", &input.Name, formscanner.StringRequired, formscanner.MinMax(4, 20)),
		newScannerFunc(r, "description", &input.Description, formscanner.String, formscanner.MinMax(4, 20)),
	}

	m := runScanners(scans[:])

	return &input, m
}

func ShowLabelCreate(t *views.Templates) http.HandlerFunc {
	return DefaultHandler(t, func(ctx context.Context, r1 *http.Request, ds *DefaultSession, r2 *views.Renderer) error {
		return r2.LabelCreateForm()
	})
}

func ShowLabelEdit(t *views.Templates, db *gorm.DB) http.HandlerFunc {
	h := func(ctx context.Context, r *http.Request, session *Session, vr *views.Renderer) error {
		var label models.Label
		if err := db.First(&label, session.ResId).Error; err != nil {
			return err
		}
		return vr.LabelEditForm(session.ResId, &label)
	}
	return ResourceHandler(t, h)
}

func ShowLabelIndex(t *views.Templates, db *gorm.DB) http.HandlerFunc {
	h := func(ctx context.Context, r *http.Request, session *DefaultSession, vr *views.Renderer) error {
		var results []models.Label
		if err := db.Find(&results).Error; err != nil {
			return err
		}
		return vr.LabelIndexPage(results)
	}
	return DefaultHandler(t, h)
}

func HandleLabelUpdate(t *views.Templates, db *gorm.DB) http.HandlerFunc {
	res := UpdateResource[models.Label]{
		parseInput: parseLabel,
		handleParseError: func(ctx context.Context, r *http.Request, session *Session, input *models.Label, fe formErrors, renderer *views.Renderer) error {
			return renderer.LabelUpdateError(session.ResId, input, fe)
		},
		handle: func(ctx context.Context, r *http.Request, s *Session, input *models.Label, vr *views.Renderer) error {

			var label models.Label
			if err := db.WithContext(ctx).First(&label, s.ResId).Error; err != nil {
				return err
			}
			updates := map[string]any{
				"Name":        input.Name,
				"Description": input.Description,
			}
			if err := db.WithContext(ctx).Model(&label).Updates(updates).Error; err != nil {
				return err
			}
			return vr.LabelUpdateOk(&label)
		},
	}
	return UpdateHandler(t, res)
}

func HandleLabelDelete(db *gorm.DB) http.HandlerFunc {
	h := func(ctx context.Context, r *http.Request, userId, resId int) error {
		var label models.Label
		if err := db.First(&label, resId).Error; err != nil {
			return err
		}
		if err := db.Delete(&label).Error; err != nil {
			return err
		}
		return nil
	}

	return DeleteHandler(h)
}

func HandleLabelCreate(t *views.Templates, db *gorm.DB) http.HandlerFunc {

	res := CreateResource[models.Label]{

		parseInput: parseLabel,
		handleParseError: func(ctx context.Context, r *http.Request, session *DefaultSession, input *models.Label, fe formErrors, vr *views.Renderer) error {
			return vr.LabelCreateError(input, fe)
		},

		handle: func(ctx context.Context, r *http.Request, session *DefaultSession, input *models.Label, vr *views.Renderer) error {

			if err := db.WithContext(r.Context()).Create(&input).Error; err != nil {
				return err
			}
			return vr.LabelCreateOk(input)
		},
	}
	return CreateHandler(t, res)
}
