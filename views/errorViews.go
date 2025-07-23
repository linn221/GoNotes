package views

import "net/http"

func ShowErrorBox(t *Templates, w http.ResponseWriter, message string, title string) {
	t.errorBoxTemplate.ExecuteTemplate(w, "", map[string]any{
		"Title":   title,
		"Message": message,
	})
}
