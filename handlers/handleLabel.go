package handlers

import (
	"context"
	"linn221/shop/formscanner"
	"linn221/shop/models"
	"linn221/shop/services"
	"linn221/shop/views"
	"net/http"
)

func parseLabel(r *http.Request) (*models.Label, services.FormErrors) {

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

func ShowLabelEdit(t *views.Templates, labelService *models.LabelService) http.HandlerFunc {
	h := func(ctx context.Context, r *http.Request, session *Session, vr *views.Renderer) error {
		label, err := labelService.Get(r.Context(), session.UserId, session.ResId)
		if err != nil {
			return err
		}
		return vr.LabelEditForm(session.ResId, label)
	}
	return ResourceHandler(t, h)
}

func ShowLabelIndex(t *views.Templates, labelService *models.LabelService) http.HandlerFunc {
	h := func(ctx context.Context, r *http.Request, session *DefaultSession, vr *views.Renderer) error {
		results, err := labelService.ListAll(r.Context(), session.UserId)
		if err != nil {
			return err
		}
		return vr.LabelIndexPage(results)
	}
	return DefaultHandler(t, h)
}

func HandleLabelUpdate(t *views.Templates, labelService *models.LabelService) http.HandlerFunc {
	handle := func(w http.ResponseWriter, r *http.Request, s *Session, input *models.Label, fe services.FormErrors, vr *views.Renderer) error {
		if len(fe) > 0 {
			return vr.LabelUpdateError(s.ResId, input, fe)
		}
		label, err := labelService.Update(r.Context(), s.UserId, s.ResId, input)
		if err != nil {
			return err
		}
		return vr.LabelUpdateOk(label)
	}
	return UpdateHandler(t, parseLabel, handle)
}

func HandleLabelDelete(labelService *models.LabelService) http.HandlerFunc {
	h := func(ctx context.Context, r *http.Request, userId, resId int) error {
		_, err := labelService.Delete(ctx, userId, resId)
		return err
	}

	return DeleteHandler(h)
}

func HandleLabelCreate(t *views.Templates, labelService *models.LabelService) http.HandlerFunc {

	handle := func(w http.ResponseWriter, r *http.Request, session *DefaultSession, input *models.Label, fe services.FormErrors, vr *views.Renderer) error {
		if len(fe) > 0 {
			w.Header().Add("HX-Retarget", "#create-form")
			w.Header().Add("HX-Reswap", "outerHTML")
			return vr.LabelCreateError(input, fe)
		}
		label, err := labelService.Create(r.Context(), session.UserId, input)
		if err != nil {
			return err
		}
		return vr.LabelCreateOk(label)
	}
	return CreateHandler(t, parseLabel, handle)
}
