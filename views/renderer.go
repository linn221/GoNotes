package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// var
//
//	(
//		loginTemplate    *template.Template
//	registerTemplate *template.Template
//	labelTemplate    *template.Template
//	testTemplate     *template.Template
//	)

type Renderer struct {
	loginTemplate *template.Template
	// registerTemplate *template.Template
	indexTemplate *template.Template

	labelTemplate *template.Template

	internalErrorTemplate  *template.Template //2d
	notFoundTemplate       *template.Template //2d
	invalidRequestTemplate *template.Template // 2d
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

func NewRenderer(baseDir string) *Renderer {
	templateDir := filepath.Join(baseDir, "views")
	parsePage := parseTemplateForPage(templateDir)
	return &Renderer{
		loginTemplate: template.Must(
			template.New("root").ParseFiles(filepath.Join(templateDir, "login.gotmpl"))),
		indexTemplate: parsePage("index.gotmpl"),

		labelTemplate: parsePage("label.gotmpl"),
	}
}

func (r *Renderer) CheckSystemError(w http.ResponseWriter, err error) {
	if err != nil {
		r.InternalServerError(w, err)
	}
}
