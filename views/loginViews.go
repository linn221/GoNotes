package views

import (
	"linn221/shop/models"
	"linn221/shop/services"
	"net/http"
)

func (r *Templates) LoginFormWithErrors(w http.ResponseWriter, input *models.User, errMap services.FormErrors) error {

	data := map[string]FormInput{
		"username": NewFormInput(input.Username, errMap["username"]),
		"password": NewFormInput(input.Password, errMap["password"]),
	}
	return r.loginTemplate.ExecuteTemplate(w, "error_view", data)
}

func (r *Templates) Index(w http.ResponseWriter) error {
	return r.indexTemplate.Execute(w, nil)
}

func (r *Templates) LoginPage(w http.ResponseWriter) error {
	return r.loginTemplate.Execute(w, nil)
}
