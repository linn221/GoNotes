package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type MyError struct {
	Error  error
	Status int
}

// will parse the request, if found errors, will write to the response
// instance, continue, internalError
func writeValidationErrors(w http.ResponseWriter, errs validator.ValidationErrors) error {
	var errors []ValidationError
	for _, err := range errs {
		errors = append(errors, ValidationError{
			Field:   err.Field(),
			Message: fmt.Sprintf("Field validation for '%s' failed on the '%s' constraint", err.Field(), err.Tag()),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(w).Encode(map[string]any{"errors": errors})
}

func respondOk(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(v)
}

func respondNotFound(w http.ResponseWriter, s string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	return json.NewEncoder(w).Encode(map[string]any{
		"error": s,
	})

}

func respondClientError(w http.ResponseWriter, s string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(w).Encode(map[string]any{
		"error": s,
	})
}

func respondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func first[T any](db *gorm.DB, shopId string, id int) (*T, *ServiceError) {
	var v T
	err := db.Where("shop_id = ?", shopId).First(&v, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &ServiceError{Code: http.StatusNotFound, Err: err}
		}
		return nil, systemErr(err)
	}

	return &v, nil
}

// will parse the request, if found errors, will write to the response
// instance, continue, internalError
func parseJson[T any](w http.ResponseWriter, r *http.Request) (*T, bool, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		finalErr := json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return nil, false, finalErr
	}

	defer r.Body.Close()
	err = validateStruct.Struct(&v)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			err := writeValidationErrors(w, ve)
			return nil, false, err
		}
		return nil, false, err
	}
	return &v, true, nil
}
