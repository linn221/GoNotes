package handlers

import (
	"linn221/shop/formscanner"
	"linn221/shop/models"
	"net/http"
)

func parseRegisterForm(r *http.Request) (*models.User, map[string]error) {
	var input models.User

	scans := [...]ScannerFunc{
		newScannerFunc(r, "email", &input.Email, formscanner.StringRequired, formscanner.MinMax(4, 20)),
		newScannerFunc(r, "username", &input.Username, formscanner.StringRequired, formscanner.MinMax(4, 20)),
		newScannerFunc(r, "password", &input.Password, formscanner.StringRequired, formscanner.MinMax(4, 100)),
	}
	m := make(map[string]error)
	for _, f := range scans {
		if field, err := f(); err != nil {
			m[field] = err
		}
	}
	return &input, m
}
