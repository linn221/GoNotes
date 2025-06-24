package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Templates struct {
	loginTemplate    *template.Template
	registerTemplate *template.Template
	indexTemplate    *template.Template

	labelTemplate          *template.Template
	noteTemplate           *template.Template
	noteCreateTemplate     *template.Template
	noteEditTemplate       *template.Template
	importNoteTemplate     *template.Template
	changePasswordTemplate *template.Template

	// internalErrorTemplate  *template.Template //2d
	// notFoundTemplate       *template.Template //2d
	// invalidRequestTemplate *template.Template // 2d
}

func parseTemplateForPage(dir string) func(filenames ...string) *template.Template {
	return func(filenames ...string) *template.Template {
		fileurls := make([]string, 0, len(filenames)+2)
		fileurls = append(fileurls, filepath.Join(dir, "layout.partial.gotmpl"),
			filepath.Join(dir, "nav.partial.gotmpl"),
		)
		for _, fn := range filenames {
			fileurls = append(fileurls, filepath.Join(dir, fn))
		}
		return template.Must(template.New("root").ParseFiles(fileurls...))
	}
}

func NewTemplates(templateDir string) *Templates {
	// templateDir := filepath.Join(baseDir, "../../views/templates")
	parsePage := parseTemplateForPage(templateDir)
	return &Templates{
		loginTemplate: template.Must(
			template.New("root").ParseFiles(filepath.Join(templateDir, "login.gotmpl"))),
		registerTemplate: template.Must(
			template.New("root").ParseFiles(filepath.Join(templateDir, "register.gotmpl"))),
		changePasswordTemplate: parsePage("change-password.gotmpl"),
		indexTemplate:          parsePage("index.gotmpl"),
		labelTemplate:          parsePage("label.gotmpl"),
		noteTemplate:           parsePage("note.gotmpl"),
		noteCreateTemplate:     parsePage("note-create.gotmpl"),
		noteEditTemplate:       parsePage("note-edit.gotmpl"),
		importNoteTemplate:     parsePage("import-notes.gotmpl"),
	}
}

type Renderer struct {
	w         http.ResponseWriter
	userId    int
	templates *Templates
}

func (t *Templates) NewRenderer(w http.ResponseWriter, userId int) *Renderer {
	return &Renderer{w: w, userId: userId, templates: t}
}
