package views

import (
	"linn221/shop/models"
	"net/http"
)

func (r *Renderer) LoginFormWithErrors(w http.ResponseWriter, input *models.User, errMap map[string]error) error {

	data := map[string]FormError{
		"username": NewFormError(input.Username, errMap["username"]),
		"password": NewFormError(input.Password, errMap["password"]),
	}
	return r.loginTemplate.ExecuteTemplate(w, "error_view", data)
}

func (r *Renderer) Index(w http.ResponseWriter) error {
	return r.indexTemplate.Execute(w, nil)
}

func (r *Renderer) Login(w http.ResponseWriter) error {
	return r.loginTemplate.Execute(w, nil)
}
