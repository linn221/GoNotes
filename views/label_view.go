package views

import (
	"linn221/shop/models"
)

func (r *Renderer) LabelCreateForm() error {
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "create_form", nil)
}

func (r *Renderer) LabelCreateError(input *models.Label, errmap map[string]error) error {

	m := map[string]FormError{
		"name":        NewFormError(input.Name, errmap["name"]),
		"description": NewFormError(input.Description, errmap["description"]),
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

func (r *Renderer) LabelUpdateError(resId int, input *models.Label, errMap map[string]error) error {

	m := map[string]any{
		"Id":          resId,
		"name":        NewFormError(input.Name, errMap["name"]),
		"description": NewFormError(input.Description, errMap["description"]),
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
