package views

import (
	"linn221/shop/models"
	"net/http"
)

type LoginInput struct {
	Username FormError
	Password FormError
}

func (r *Renderer) LoginFormWithErrors(w http.ResponseWriter, input *models.User, errMap map[string]error) error {

	page := newFormErrorPage("login")
	page.Input = map[string]FormError{
		"username": NewFormError(input.Username, errMap["username"]),
		"password": NewFormError(input.Password, errMap["password"]),
	}
	return r.loginTemplate.Execute(w, page)
}

func (r *Renderer) Index(w http.ResponseWriter) error {
	return r.indexTemplate.Execute(w, nil)
}

func (r *Renderer) Login(w http.ResponseWriter) error {
	page := newCreatePage("Login")
	return r.loginTemplate.Execute(w, page)
}
