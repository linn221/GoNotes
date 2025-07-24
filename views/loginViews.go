package views

import (
	"linn221/shop/models"
	"linn221/shop/services"
	"net/http"
)

func (r *Templates) Index(w http.ResponseWriter) error {
	return r.indexTemplate.Execute(w, nil)
}

func (r *Templates) LoginPage(w http.ResponseWriter) error {
	return r.loginTemplate.Execute(w, nil)
}

func (r *Templates) RegisterFormWithErrors(w http.ResponseWriter, input *models.User, errMap services.FormErrors) error {

	data := map[string]FormInput{
		"username": NewFormInput(input.Username, errMap["username"]),
		"email":    NewFormInput(input.Email, errMap["email"]),
		"password": NewFormInput(input.Password, errMap["password"]),
	}
	return r.registerTemplate.ExecuteTemplate(w, "error_view", data)
}
func (r *Templates) RegisterPage(w http.ResponseWriter) error {
	return r.registerTemplate.Execute(w, nil)
}

func (t *Templates) ChangePasswordPage(w http.ResponseWriter) error {
	return t.changePasswordTemplate.Execute(w, nil)
}
func (t *Templates) ChangePasswordWithError(w http.ResponseWriter, data map[string]FormInput) error {
	return t.changePasswordTemplate.ExecuteTemplate(w, "error_view", data)
}
