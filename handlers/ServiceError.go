package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ServiceError struct {
	Err  error
	Code int
}

func (se *ServiceError) Respond(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(se.Code)
	return json.NewEncoder(w).Encode(map[string]any{
		"error": se.Err.Error(),
	})
}
func systemErr(err error) *ServiceError {
	return &ServiceError{
		Err:  err,
		Code: http.StatusInternalServerError,
	}
}
func systemErrString(s string, args ...error) *ServiceError {
	if len(args) > 0 {
		err := args[0]
		return systemErr(errors.New(s + ": " + err.Error()))
	}
	return systemErr(errors.New(s))
}
func clientErr(s string) *ServiceError {
	return &ServiceError{
		Err:  errors.New(s),
		Code: http.StatusBadRequest,
	}
}
func ServiceErrorHandler(h func(http.ResponseWriter, *http.Request) *ServiceError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serr := h(w, r)
		if serr != nil {
			if err := serr.Respond(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
