package handlers

import (
	"linn221/shop/formscanner"
	"linn221/shop/models"
	"linn221/shop/services"
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
