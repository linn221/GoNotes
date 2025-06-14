package views

import (
	"linn221/shop/models"
	"linn221/shop/services"
)

func (r *Renderer) LabelCreateForm() error {
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "create_form", nil)
}

func (r *Renderer) LabelCreateError(input *models.Label, errmap services.FormErrors) error {

	m := map[string]FormInput{
		"name":        NewFormInput(input.Name, errmap["name"]),
		"description": NewFormInput(input.Description, errmap["description"]),
	}
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "create_form2", m)
}

func (r *Renderer) LabelCreateOk(label *models.Label) error {
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "create_success", map[string]any{
		"Res": label,
	})
}

func (r *Renderer) LabelUpdateOk(label *models.Label) error {
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "edit_success", map[string]any{
		"Res": label,
	})
}

func (r *Renderer) LabelEditForm(resId int, input *models.Label) error {
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "edit_form", map[string]any{
		"Res": input,
		"Id":  resId,
	})
}

func (r *Renderer) LabelUpdateError(resId int, input *models.Label, errMap services.FormErrors) error {

	m := map[string]any{
		"Id":          resId,
		"name":        NewFormInput(input.Name, errMap["name"]),
		"description": NewFormInput(input.Description, errMap["description"]),
	}
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "edit_form2", m)
}

func (r *Renderer) LabelIndexPage(labels []models.Label) error {
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "root", map[string]any{
		"ResList": labels,
	})
}

// func (r *Renderer) LabelDetails(w http.ResponseWriter, label *models.Label) error {
// 	page := newIndexPage("label details")
// 	page.Res = label
// 	return r.labelDetailsTemplate.Execute(w, page)
// }
