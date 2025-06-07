package handlers

import (
	"linn221/shop/formscanner"
	"net/http"
)

// type ValidationError struct {
// 	Field   string `json:"field"`
// 	Message string `json:"message"`
// }

// type MyError struct {
// 	Error  error
// 	Status int
// }

// func respondOk(w http.ResponseWriter, v any) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	return json.NewEncoder(w).Encode(v)
// }

// func respondNotFound(w http.ResponseWriter, s string) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusNotFound)
// 	return json.NewEncoder(w).Encode(map[string]any{
// 		"error": s,
// 	})

// }

// func respondClientError(w http.ResponseWriter, s string) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusBadRequest)
// 	return json.NewEncoder(w).Encode(map[string]any{
// 		"error": s,
// 	})
// }

// func respondNoContent(w http.ResponseWriter) {
// 	w.WriteHeader(http.StatusNoContent)
// }

// func first[T any](db *gorm.DB, shopId string, id int) (*T, *ServiceError) {
// 	var v T
// 	err := db.Where("shop_id = ?", shopId).First(&v, id).Error
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, &ServiceError{Code: http.StatusNotFound, Err: err}
// 		}
// 		return nil, systemErr(err)
// 	}

// 	return &v, nil
// }

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

func createdResponse(w http.ResponseWriter) {
	w.Header().Add("HX-Trigger", "content-created")
	w.WriteHeader(http.StatusNoContent)
}

func updatedResponse(w http.ResponseWriter) {
	w.Header().Add("HX-Trigger", "content-updated")
	w.WriteHeader(http.StatusNoContent)
}

func finalErrHandle(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
