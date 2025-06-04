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
	"gorm.io/gorm"
)

func parseLoginForm(r *http.Request) (*models.User, map[string]error) {
	var input models.User
	scans := [...]ScanFormFunc{
		newScanner(r, "username", &input.Username, formscanner.StringRequired, formscanner.MinMax(4, 20)),
		newScanner(r, "password", &input.Password, formscanner.StringRequired, formscanner.MinMax(3, 20)),
	}

	m := runScanners(scans[:])
	return &input, m
}

func HandleLogin(vr *views.Renderer,
	db *gorm.DB,
	cache services.CacheService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			finalErrHandle(w,
				vr.Login(w),
			)
			return
		} else if r.Method == http.MethodPost {

			input, errMap := parseLoginForm(r)
			if len(errMap) > 0 {
				finalErrHandle(w,
					vr.LoginFormWithErrors(w, input, errMap),
				)
				return
			}

			var user models.User
			if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
				finalErrHandle(w,
					vr.LoginFormWithErrors(w, input, map[string]error{
						"username": errors.New("invalid username/password"),
						"password": errors.New("invalid username/password")}),
				)
				return
			}

			if err := utils.ComparePassword(user.Password, input.Password); err != nil {
				finalErrHandle(w,
					vr.LoginFormWithErrors(w, input, map[string]error{
						"username": errors.New("invalid username/password"),
						"password": errors.New("invalid username/password"),
					}),
				)
				return
			}

			token, err := newSessionToken(cache, user.Id)
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
			// w.Header().Set("HX-Redirect", "/")
			// w.WriteHeader(http.StatusOK)
			// utils.HxRedirect(w, "/users/labels")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		} else {
			finalErrHandle(w, errors.New("invalid http method"))
		}
	}
}

func newSessionToken(cache services.CacheService, userId int) (string, error) {
	tokenString := uuid.NewString()
	if err := cache.SetValue(fmt.Sprintf("Token:%s", tokenString), fmt.Sprint(userId), time.Hour*127); err != nil {
		return "", err
	}

	return tokenString, nil
}

func HandleLogout(cache services.CacheService, vr *views.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		//2d write helper for removing token from redis
		if err := cache.RemoveKey(fmt.Sprintf("Token:%s", token)); err != nil {
			finalErrHandle(w, err)
			return
		}
		removeTokenCookies(w)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}
}

func removeTokenCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Unix(0, 0), // Set to past
		MaxAge:  -1,              // Also ensures deletion
		Path:    "/",
		Domain:  "",
	})
}
