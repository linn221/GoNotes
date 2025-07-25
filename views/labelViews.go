package views

import (
	"linn221/shop/models"
)

func (r *Renderer) ShowLabelCreate() error {
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "create_form", nil)
}

func (r *Renderer) HandleLabelCreate(label *models.Label) error {
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "create_success", ResourceData{Res: label})
}

func (r *Renderer) HandleLabelUpdate(label *models.Label) error {
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "edit_success", ResourceData{Res: label})
}

func (r *Renderer) ShowLabelEdit(resId int, label *models.Label) error {
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "edit_form", ResourceData{Res: label})
}

func (r *Renderer) HandleLabelToggleActive(label *models.Label) error {
	return r.templates.labelTemplate.ExecuteTemplate(r.w, "toggle_button", label)
}

func (r *Renderer) RenderLabelIndex(labels []models.Label) error {
	return r.templates.labelTemplate.Execute(r.w, Page{
		Nav:       NavLabels,
		PageTitle: "Labels",
		ResList:   labels,
	})
}

// func (r *Renderer) LabelDetails(w http.ResponseWriter, label *models.Label) error {
// 	page := newIndexPage("label details")
// 	page.Res = label
// 	return r.labelDetailsTemplate.Execute(w, page)
// }
