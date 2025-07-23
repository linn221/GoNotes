package views

import (
	"linn221/shop/services"
	"net/http"
)

func ShowFormError(t *Templates, w http.ResponseWriter, fe services.FormErrors) error {
	w.Header().Add("HX-Reswap", "outerHTML")
	w.Header().Add("HX-Retarget", "#flash") // oob swap does not work with after-swap event so retargeting it
	return t.errorBoxTemplate.ExecuteTemplate(w, "form_error", map[string]any{
		"ErrorMap": fe,
	})
}

func ShowDefaultError(t *Templates, w http.ResponseWriter, err error) error {
	w.Header().Add("HX-Reswap", "outerHTML")
	w.Header().Add("HX-Retarget", "#flash") // oob swap does not work with after-swap event so retargeting it
	return t.errorBoxTemplate.ExecuteTemplate(w, "error", map[string]any{
		"Title":   "Error",
		"Message": err.Error(),
	})
}

func ShowInternalError(t *Templates, w http.ResponseWriter, err error) error {
	w.Header().Add("HX-Reswap", "outerHTML")
	w.Header().Add("HX-Retarget", "#flash")
	return t.errorBoxTemplate.ExecuteTemplate(w, "error", map[string]any{
		"Title":   "Internal Error",
		"Message": err.Error(),
	})
}
