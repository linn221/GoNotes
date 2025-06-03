package handlers

import (
	"linn221/shop/models"
	"linn221/shop/services"
	"net/http"

	"gorm.io/gorm"
)

type NewUnit struct {
	Name        inputString     `json:"name" validate:"required,min=2,max=100"`
	Symbol      inputString     `json:"symbol" validate:"required,min=1,max=10"`
	Description *optionalString `json:"description" validate:"omitempty,max=500"`
}

func (input *NewUnit) validate(db *gorm.DB, shopId string, id int) *ServiceError {

	shopFilter := NewShopFilter(shopId)
	if err := Validate(db,
		NewExistsRule("units", id, "unit not found", shopFilter).When(id > 0),
		NewUniqueRule("units", "name", input.Name, id, "duplicate name", shopFilter),
		NewUniqueRule("units", "symbol", input.Symbol, id, "duplicate symbol", shopFilter),
	); err != nil {
		return err
	}
	return nil
}

func HandleUnitCreate(db *gorm.DB,
	cleanListingCache services.CleanListingCache,
) http.HandlerFunc {
	return CreateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewUnit) error {

		if errs := input.validate(db.WithContext(r.Context()), session.ShopId, 0); errs != nil {
			return errs.Respond(w)
		}
		unit := models.Unit{
			Name:        input.Name.String(),
			Symbol:      input.Symbol.String(),
			Description: input.Description.StringPtr(),
		}
		unit.ShopId = session.ShopId
		err := db.WithContext(r.Context()).Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&unit).Error; err != nil {
				return err
			}

			if err := cleanListingCache(session.ShopId); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)
		return nil
	})
}

func HandleUnitUpdate(db *gorm.DB,
	cleanCache func(db *gorm.DB, shopId string, id int) error,
) http.HandlerFunc {

	return UpdateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewUnit) error {

		ctx := r.Context()
		if errs := input.validate(db.WithContext(ctx), session.ShopId, session.ResId); errs != nil {
			return errs.Respond(w)
		}
		updates := map[string]any{
			"Name":   input.Name.String(),
			"Symbol": input.Symbol.String(),
		}
		if input.Description.IsPresent() {
			updates["Description"] = input.Description.String()
		}
		unit, errs := first[models.Unit](db.WithContext(ctx), session.ShopId, session.ResId)
		if errs != nil {
			return errs.Respond(w)
		}

		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&unit).Updates(updates).Error; err != nil {
				return err
			}

			if err := cleanCache(db.WithContext(ctx), session.ShopId, session.ResId); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})
}

func HandleUnitDelete(db *gorm.DB,
	cleanCache func(db *gorm.DB, shopId string, id int) error,
) http.HandlerFunc {
	return DeleteHandler(func(w http.ResponseWriter, r *http.Request, session Session) error {
		ctx := r.Context()
		unit, errs := first[models.Unit](db.WithContext(ctx), session.ShopId, session.ResId)
		if errs != nil {
			return errs.Respond(w)
		}

		if errs := Validate(db.WithContext(ctx),
			NewNoResultRule("units", "unit has been used in items", NewFilter("unit_id = ?", session.ResId)),
		); errs != nil {
			return errs.Respond(w)
		}
		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&unit).Error; err != nil {
				return err
			}
			if err := cleanCache(db.WithContext(ctx), session.ShopId, session.ResId); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})
}
