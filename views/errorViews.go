package views

import "net/http"

func (r *Renderer) InternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	r.internalErrorTemplate.Execute(w, err.Error())
}

func (r *Renderer) NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	r.notFoundTemplate.Execute(w, nil)
}

func (r *Renderer) InvalidRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	r.invalidRequestTemplate.Execute(w, err)
}
