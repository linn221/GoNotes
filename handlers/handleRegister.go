package handlers

import (
	"linn221/shop/formscanner"
	"linn221/shop/models"
	"linn221/shop/services"
	"linn221/shop/views"
	"net/http"
)

func parseRegisterForm(r *http.Request) (*models.User, services.FormErrors) {
	var input models.User

	scans := [...]ScannerFunc{
		newScannerFunc(r, "email", &input.Email, formscanner.StringRequired, formscanner.MinMax(4, 20)),
		newScannerFunc(r, "username", &input.Username, formscanner.StringRequired, formscanner.MinMax(4, 20)),
		newScannerFunc(r, "password", &input.Password, formscanner.StringRequired, formscanner.MinMax(4, 100)),
	}
	m := make(services.FormErrors)
	for _, f := range scans {
		if field, err := f(); err != nil {
			m[field] = err
		}
	}
	return &input, m
}

func HandleRegister(t *views.Templates, userSevice *models.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			finalErrHandle(w, t.RegisterPage(w))
			return
		} else if r.Method == http.MethodPost {
			input, errmap := parseRegisterForm(r)
			if len(errmap) > 0 {
				finalErrHandle(w, t.RegisterFormWithErrors(w, input, errmap))
				return
			}
			_, err := userSevice.Register(r.Context(), input)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			htmxRedirect(w, "/login")
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
