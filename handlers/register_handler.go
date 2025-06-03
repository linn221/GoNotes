package handlers

import (
	"linn221/shop/formscanner"
	"linn221/shop/models"
	"net/http"
)

func parseRegisterForm(r *http.Request) (*models.User, map[string]error) {
	var input models.User

	scans := [...]ScanFormFunc{
		newScanner(r, "email", &input.Email, formscanner.StringRequired, formscanner.MinMax(4, 20)),
		newScanner(r, "username", &input.Username, formscanner.StringRequired, formscanner.MinMax(4, 20)),
		newScanner(r, "password", &input.Password, formscanner.StringRequired, formscanner.MinMax(4, 100)),
	}
	m := make(map[string]error)
	for _, f := range scans {
		if field, err := f(); err != nil {
			m[field] = err
		}
	}
	return &input, m
}

// func HandleRegister(db *gorm.DB,
// 	vr *views.Renderer,
// ) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodGet {
// 			return tmpl.Execute(w, nil)
// 		} else if r.Method == http.MethodPost {
// 			input, errMap := parseRegisterForm(r)
// 			if len(errMap) > 0 {
// 				return tmpl.ExecuteTemplate(w, "form", map[string]any{
// 					"HasErrors": true,
// 					"Errors":    errMap,
// 					"Data":      input,
// 				})
// 			}

// 			hashedPassword, err := utils.HashPassword(input.Password)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusBadRequest)
// 			}
// 			input.Password = string(hashedPassword)
// 			if err = db.Create(&input).Error; err != nil {
// 				return err
// 			}

// 			// w.Header().Set("HX-Redirect", "/login")
// 			w.WriteHeader(http.StatusOK)
// 			return nil
// 		}
// 		http.Error(w, "invalid http method", http.StatusBadRequest)
// 		return nil
// 	}
// }
