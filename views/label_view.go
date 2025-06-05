package views

import (
	"linn221/shop/models"
	"net/http"
)

func (r *Renderer) LabelCreateForm(w http.ResponseWriter, userId int) error {
	return r.labelTemplate.Execute(w, nil)
}

func (r *Renderer) LabelCreateError(w http.ResponseWriter, userId int, input *models.Label, errmap map[string]error) error {

	m := map[string]FormError{
		"name":        NewFormError(input.Name, errmap["name"]),
		"description": NewFormError(input.Description, errmap["description"]),
	}
	return r.labelTemplate.ExecuteTemplate(w, "create_form2", m)
}

func (r *Renderer) LabelCreateOk(w http.ResponseWriter, label *models.Label) error {
	return r.labelTemplate.ExecuteTemplate(w, "create_success", map[string]any{
		"Res": label,
	})
}

func (r *Renderer) LabelUpdateOk(w http.ResponseWriter, label *models.Label) error {
	return r.labelTemplate.ExecuteTemplate(w, "edit_success", map[string]any{
		"Res": label,
	})
}

func (r *Renderer) LabelEditForm(w http.ResponseWriter, userId int, resId int, input *models.Label) error {
	return r.labelTemplate.ExecuteTemplate(w, "edit_form", map[string]any{
		"Res": input,
		"Id":  resId,
	})
}
func (r *Renderer) LabelUpdateError(w http.ResponseWriter, userId int, resId int, input *models.Label, errMap map[string]error) error {

	m := map[string]any{
		"Id":          resId,
		"name":        NewFormError(input.Name, errMap["name"]),
		"description": NewFormError(input.Description, errMap["description"]),
	}
	return r.labelTemplate.ExecuteTemplate(w, "edit_form2", m)
}

func (r *Renderer) LabelIndexPage(w http.ResponseWriter, labels []models.Label, loadFullPage bool) error {
	templateName := "main"
	if loadFullPage {
		templateName = "root"
	}
	return r.labelTemplate.ExecuteTemplate(w, templateName, map[string]any{
		"ResList": labels,
	})
}

// func (r *Renderer) LabelDetails(w http.ResponseWriter, label *models.Label) error {
// 	page := newIndexPage("label details")
// 	page.Res = label
// 	return r.labelDetailsTemplate.Execute(w, page)
// }
