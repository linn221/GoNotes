package views

import "net/http"

func (t *Templates) AccountPage(w http.ResponseWriter) error {
	return t.accountTemplate.ExecuteTemplate(w, "root", H{
		"Nav":       NavAccount,
		"PageTitle": "Account actions",
	})
}
