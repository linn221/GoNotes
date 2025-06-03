package handlers

import (
	"linn221/shop/models"
	"linn221/shop/services"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type NewItem struct {
	Name          inputString      `json:"name" validate:"required,min=2,max=100"`
	Description   *optionalString  `json:"description" validate:"omitempty,max=500"`
	CategoryId    int              `json:"category_id" validate:"required,number,gte=1"`
	UnitId        int              `json:"unit_id" validate:"required,number,gte=1"`
	SalesPrice    *decimal.Decimal `json:"sales_price" validate:"required,number,gte=1"`
	PurchasePrice *decimal.Decimal `json:"purchase_price" validate:"required,number,gte=1"`
}

func (input *NewItem) validate(db *gorm.DB, shopId string, id int) *ServiceError {
	shopFilter := NewShopFilter(shopId)
	if err := Validate(db,
		NewExistsRule("items", id, "item not found", shopFilter).When(id > 0),
		NewUniqueRule("items", "name", input.Name, id, "duplicate item name", shopFilter),
		NewExistsRule("units", input.UnitId, "unit not found", shopFilter),
		NewExistsRule("categories", input.CategoryId, "category not found", shopFilter),
	); err != nil {
		return err
	}
	return nil
}

func HandleItemCreate(db *gorm.DB,
	cleanListingCache services.CleanListingCache,
) http.HandlerFunc {
	return CreateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewItem) error {

		ctx := r.Context()
		if errs := input.validate(db.WithContext(ctx), session.ShopId, 0); errs != nil {
			return errs.Respond(w)
		}
		item := models.Item{
			Name:          input.Name.String(),
			Description:   input.Description.StringPtr(),
			CategoryId:    input.CategoryId,
			UnitId:        input.UnitId,
			SalesPrice:    *input.SalesPrice,
			PurchasePrice: *input.PurchasePrice,
		}
		item.ShopId = session.ShopId

		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&item).Error; err != nil {
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

func HandleItemUpdate(db *gorm.DB,
	cleanInstanceCache services.CleanInstanceCache,
	cleanListingCache services.CleanListingCache,
) http.HandlerFunc {
	return UpdateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewItem) error {

		ctx := r.Context()
		if errs := input.validate(db.WithContext(ctx), session.ShopId, session.ResId); errs != nil {
			return errs.Respond(w)
		}
		item, errs := first[models.Item](db.WithContext(ctx), session.ShopId, session.ResId)
		if errs != nil {
			return errs.Respond(w)
		}

		updates := map[string]any{
			"Name":          input.Name,
			"CategoryId":    input.CategoryId,
			"UnitId":        input.UnitId,
			"PurchasePrice": input.PurchasePrice,
			"SalesPrice":    input.SalesPrice,
		}
		if input.Description.IsPresent() {
			updates["Description"] = input.Description.String()
		}

		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&item).Updates(updates).Error; err != nil {
				return err
			}
			if err := cleanInstanceCache(session.ResId); err != nil {
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

		respondNoContent(w)
		return nil
	})
}

func HandleItemDelete(db *gorm.DB,
	cleanInstanceCache services.CleanInstanceCache,
	cleanListingCache services.CleanListingCache,
) http.HandlerFunc {
	return DeleteHandler(func(w http.ResponseWriter, r *http.Request, session Session) error {

		ctx := r.Context()
		item, errs := first[models.Item](db.WithContext(ctx), session.ShopId, session.ResId)
		if errs != nil {
			return errs.Respond(w)
		}

		// valiate
		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&item).Error; err != nil {
				return err
			}
			if err := cleanInstanceCache(session.ResId); err != nil {
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

		respondNoContent(w)
		return nil
	})
}

func parseItemSearch(r *http.Request) (*models.ItemSearch, error) {
	var search models.ItemSearch
	var err error
	if s := r.URL.Query().Get("search"); s != "" {
		search.Search = s
	}
	if s := r.URL.Query().Get("category_id"); s != "" {
		search.CategoryId, err = strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
	}
	if s := r.URL.Query().Get("unit_id"); s != "" {
		search.UnitId, err = strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
	}

	if min := r.URL.Query().Get("sales_price_min"); min != "" {
		salesPriceMin, err := decimal.NewFromString(min)
		if err != nil {
			return nil, err
		}
		search.SalesPriceMin = &salesPriceMin
		max := r.URL.Query().Get("sales_price_max")
		if max != "" {
			salesPriceMax, err := decimal.NewFromString(max)
			if err != nil {
				return nil, err
			}
			search.SalesPriceMax = &salesPriceMax
		}
	}

	if min := r.URL.Query().Get("purchase_price_min"); min != "" {
		purchasePriceMin, err := decimal.NewFromString(min)
		if err != nil {
			return nil, err
		}
		search.PurchasePriceMin = &purchasePriceMin
		max := r.URL.Query().Get("purchase_price_max")
		if max != "" {
			purchasePriceMax, err := decimal.NewFromString(max)
			if err != nil {
				return nil, err
			}
			search.PurchasePriceMax = &purchasePriceMax
		}
	}
	return &search, nil
}

func HandleItemIndex(db *gorm.DB) http.HandlerFunc {
	return DefaultHandler(func(w http.ResponseWriter, r *http.Request, session *DefaultSession) error {
		search, err := parseItemSearch(r)
		if err != nil {
			return respondClientError(w, err.Error())
		}

		//2d cache results
		results, err := models.SearchItems(db.WithContext(r.Context()), session.ShopId, search)
		if err != nil {
			return err
		}
		respondOk(w, results)
		return nil
	})
}
