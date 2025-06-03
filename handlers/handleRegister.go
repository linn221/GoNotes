package handlers

import (
	"encoding/json"
	"fmt"
	"linn221/shop/models"
	"linn221/shop/utils"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewShop struct {
	Name     string `json:"name" validate:"required,min=4,max=255"`
	Email    string `json:"email" validate:"required,email,min=5,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
	PhoneNo  string `json:"phone_no" validate:"required,min=2"`
}

func Register(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var input NewShop
		// defer r.Body.Close()
		// err := json.NewDecoder(r.Body).Decode(&input)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }
		input, ok, finalErr := parseJson[NewShop](w, r)
		if !ok {
			if finalErr != nil {
				http.Error(w, finalErr.Error(), http.StatusInternalServerError)
			}
			return
		}

		ctx := r.Context()
		tx := db.WithContext(ctx).Begin()
		defer tx.Rollback()

		shopId := uuid.NewString()
		shop := models.Shop{
			Id:      shopId,
			Name:    input.Name,
			LogoUrl: "",
			Email:   input.Email,
			PhoneNo: input.PhoneNo,
		}
		if err := tx.Create(&shop).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		i := rand.Intn(100000)
		username := fmt.Sprintf("owner%d", i)
		password, err := utils.HashPassword(input.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user := models.User{
			Username: username,
			Email:    input.Email,
			PhoneNo:  input.PhoneNo,
			Password: string(password),
		}
		user.ShopId = shopId
		if err := tx.Create(&user).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tx.Commit().Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
