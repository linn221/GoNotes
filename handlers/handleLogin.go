package handlers

import (
	"errors"
	"fmt"
	"linn221/shop/formscanner"
	"linn221/shop/models"
	"linn221/shop/services"
	"linn221/shop/utils"
	"linn221/shop/views"

	"net/http"
	"time"

	"github.com/google/uuid"
)

func parseLoginForm(r *http.Request) (*models.User, string, services.FormErrors) {
	var input models.User
	var timezone string
	scans := [...]ScannerFunc{
		newScannerFunc(r, "username", &input.Username, formscanner.StringRequired, formscanner.MinMax(4, 20)),
		newScannerFunc(r, "password", &input.Password, formscanner.StringRequired, formscanner.MinMax(3, 20)),
		newScannerFunc(r, "timezone", &timezone, formscanner.StringRequired),
	}

	m := runScanners(scans[:])
	return &input, timezone, m
}

func HandleLogin(vr *views.Templates,
	userService *models.UserService,
	cache services.CacheService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			finalErrHandle(w,
				vr.LoginPage(w),
			)
			return
		} else if r.Method == http.MethodPost {

			input, timezone, errMap := parseLoginForm(r)
			if len(errMap) > 0 {
				finalErrHandle(w,
					vr.LoginFormWithErrors(w, input, errMap),
				)
				return
			}

			user, err := userService.Login(r.Context(), input.Username, input.Password)
			if err != nil {
				finalErrHandle(w,
					vr.LoginFormWithErrors(w, input, services.FormErrors{
						"username": errors.New("invalid username/password"),
						"password": errors.New("invalid username/password")}),
				)
				return
			}

			token, err := newSessionToken(cache, user.Id, timezone)
			if err != nil {
				finalErrHandle(w, err)
				return
			}

			// set cookies
			http.SetCookie(w, &http.Cookie{
				Name:   "token",
				Value:  token,
				MaxAge: 0,
				Path:   "/", Domain: "",
				Secure: false, HttpOnly: true,
			})
			htmxRedirect(w, "/labels")
		} else {
			finalErrHandle(w, errors.New("invalid http method"))
		}
	}
}

func newSessionToken(cache services.CacheService, userId int, timezone string) (string, error) {
	tokenString := uuid.NewString()
	if err := cache.SetH(fmt.Sprintf("Token:%s", tokenString), map[string]any{
		"timezone": timezone,
		"userId":   userId,
	}, time.Hour*127); err != nil {
		return "", err
	}

	return tokenString, nil
}

func HandleLogout(cache services.CacheService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logout(cache, w, r)
	}
}

// handle logout
func logout(cache services.CacheService, w http.ResponseWriter, r *http.Request) {

	cookies, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		finalErrHandle(w, err)
		return
	}
	token := cookies.Value
	if err := cache.RemoveKey(fmt.Sprintf("Token:%s", token)); err != nil {
		finalErrHandle(w, err)
		return
	}
	utils.RemoveTokenCookies(w)
	htmxRedirect(w, "/login")
}
