package handlers

import (
	"linn221/shop/models"
	"linn221/shop/utils"
	"net/http"

	"gorm.io/gorm"
)

type NewPassword struct {
	OldPassword string `json:"old_password" validate:"required,max=255"`
	NewPassword string `json:"new_password" validate:"required,max=255,min=8"`
}

func ChangePassword(db *gorm.DB) http.HandlerFunc {
	return InputHandler(func(w http.ResponseWriter, r *http.Request, session *DefaultSession, input *NewPassword) error {
		var user models.User
		if err := db.WithContext(r.Context()).First(&user, session.UserId).Error; err != nil {
			return err
		}

		if err := utils.ComparePassword(user.Password, input.OldPassword); err != nil {
			return respondClientError(w, "passwords do not match")
		}
		hashed, err := utils.HashPassword(input.NewPassword)
		if err != nil {
			return err
		}
		if err := db.Model(&user).UpdateColumn("password", hashed).Error; err != nil {
			return err
		}
		respondNoContent(w)
		return nil
	})
}

type NewUserEdit struct {
	Username inputString     `json:"username" validate:"required,min=3,max=100"`
	Email    inputString     `json:"email" validate:"required,email,min=4,max=100"`
	PhoneNo  *optionalString `json:"phone_no" validate:"omitempty,min=5,max=20"`
}

func UpdateUserInfo(db *gorm.DB) http.HandlerFunc {
	return InputHandler(func(w http.ResponseWriter, r *http.Request, session *DefaultSession, input *NewUserEdit) error {
		var user models.User
		if err := db.WithContext(r.Context()).First(&user, session.UserId).Error; err != nil {
			return err
		}

		updates := map[string]any{
			"Username": input.Username,
			"Email":    input.Email,
		}
		if input.PhoneNo.IsPresent() {
			updates["PhoneNo"] = input.PhoneNo.String()
		}
		shopFilter := NewShopFilter(session.ShopId)
		if errs := Validate(db,
			NewUniqueRule("users", "username", input.Username, session.UserId, "duplicate username", shopFilter),
			NewUniqueRule("users", "email", input.Email, session.UserId, "duplicate email", shopFilter),
			NewUniqueRule("users", "phone_no", input.PhoneNo, session.UserId, "duplicate phone number", shopFilter).When(input.PhoneNo.IsPresent()),
		); errs != nil {
			return errs.Respond(w)
		}
		if err := db.Model(&user).Updates(updates).Error; err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})
}

type MyInfo struct {
	Id       int    `json:"id"`
	ShopId   string `json:"shop_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	PhoneNo  string `json:"phone_no"`
}

func Me(db *gorm.DB) http.HandlerFunc { // only place where shop_id should be included

	return DefaultHandler(func(w http.ResponseWriter, r *http.Request, session *DefaultSession) error {
		var user models.User
		if err := db.WithContext(r.Context()).First(&user, session.UserId).Error; err != nil {
			return err
		}
		return respondOk(w, MyInfo{
			Id:       session.UserId,
			ShopId:   session.ShopId,
			Username: user.Username,
			Email:    user.Email,
			PhoneNo:  user.PhoneNo,
		})
	})
}

// func CreateItemCategory(db *gorm.DB, cache services.CacheService, categoryService services.CategoryCruder, validate *validator.Validate) http.HandlerFunc {
// 	return inputHandler(func(w http.ResponseWriter, r *http.Request, input *models.Category, userId int, shopId string) error {
// 		cat, err := categoryService.CreateCategory(input, db, cache)
// 		if err != nil {
// 			return err
// 		}

// 		return respondOk(w, cat)
// 	}, validate)
// }

// func inputHandler[T any](f func(w http.ResponseWriter, r *http.Request, input *T, userId int, shopId string) error, validate *validator.Validate) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		userId, shopId, err := myctx.GetIdsFromContext(ctx)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		input, ok, err := parseJson[T](w, r, validate)
// 		if !ok {
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 			}
// 			return
// 		}

// 		err = f(w, r, input, userId, shopId)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// }
