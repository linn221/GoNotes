package handlers

import (
	"linn221/shop/formscanner"
	"net/http"
)

type ScannerFunc func() (string, error)

func newScannerFunc[T any](r *http.Request, name string, ptr *T, scanFunc func(*http.Request, string) (T, bool, error), validateFuncs ...formscanner.ValidateFunc[T]) ScannerFunc {
	return func() (string, error) {
		err := formscanner.Scan(r, name, ptr, scanFunc, validateFuncs...)
		if err != nil {
			return name, err
		}
		return "", nil
	}
}

func runScanners(xs []ScannerFunc) formErrors {
	m := make(map[string]error)
	for _, x := range xs {
		if inputName, err := x(); err != nil {
			m[inputName] = err
		}
	}

	return m
}
func finalErrHandle(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
