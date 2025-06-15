package handlers

import (
	"linn221/shop/formscanner"
	"linn221/shop/models"
	"linn221/shop/myctx"
	"linn221/shop/services"
	"linn221/shop/views"
	"net/http"
)

type PasswordInfo struct {
	NewPassword string
	OldPassword string
}

func parseChangePasswordForm(r *http.Request) (*PasswordInfo, services.FormErrors) {

	var input PasswordInfo
	scans := [...]ScannerFunc{
		newScannerFunc(r, "old_password", &input.OldPassword, formscanner.StringRequired, formscanner.MinMax(8, 20)),
		newScannerFunc(r, "new_password", &input.NewPassword, formscanner.StringRequired, formscanner.MinMax(8, 20)),
	}

	m := runScanners(scans[:])
	return &input, m
}

func HandleChangePassword(t *views.Templates, userService *models.UserService, cache services.CacheService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			finalErrHandle(w, t.ChangePasswordPage(w))
			return
		} else if r.Method == http.MethodPost {
			input, fe := parseChangePasswordForm(r)
			if len(fe) > 0 {
				finalErrHandle(w,
					t.ChangePasswordWithError(w, map[string]views.FormInput{
						"old_password": views.NewFormInput(input.OldPassword, fe["old_password"]),
						"new_password": views.NewFormInput(input.NewPassword, fe["new_password"]),
					}),
				)
				return
			}
			userId, err := myctx.GetUserId(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = userService.ChangePassword(r.Context(), userId, input.OldPassword, input.NewPassword)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			logout(cache, w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}
