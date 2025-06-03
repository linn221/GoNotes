package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Item struct {
	Id            int     `gorm:"primaryKey"`
	Name          string  `gorm:"index;not null"`
	Description   *string `gorm:"default:null"`
	CategoryId    int     `gorm:"index;not null"`
	SalesPrice    decimal.Decimal
	PurchasePrice decimal.Decimal
	UnitId        int `gorm:"index;not null"`
	Category      Category
	Unit          Unit
	HasShopId
	HasIsActive
}

type ItemResource struct {
	Id            int             `json:"id"`
	Name          string          `json:"name"`
	SalesPrice    decimal.Decimal `json:"sales_price"`
	PurchasePrice decimal.Decimal `json:"purchase_price"`
	CategoryName  string          `json:"category_name"`
	UnitName      string          `json:"unit_name"`
	UnitSymbol    string          `json:"unit_symbol"`
}

type ItemDetailResource struct {
	Id            int             `json:"id"`
	Name          string          `json:"name"`
	Description   *string         `json:"description"`
	SalesPrice    decimal.Decimal `json:"sales_price"`
	PurchasePrice decimal.Decimal `json:"purchase_price"`
	Category      struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		IsActive bool   `json:"is_active"`
	} `json:"category"`
	Unit struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
		IsActive bool   `json:"is_active"`
	} `json:"unit"`
	HasIsActive
	HasShopId
}

func FetchItemResources(db *gorm.DB, shopId string) ([]ItemResource, error) {
	var items []Item
	if err := db.Preload("Category").Preload("Unit").Where("shop_id = ? AND is_active = 1", shopId).Find(&items).Error; err != nil {
		return nil, err
	}
	results := make([]ItemResource, 0, len(items))
	for _, item := range items {
		results = append(results, ItemResource{
			Id:            item.Id,
			Name:          item.Name,
			SalesPrice:    item.SalesPrice,
			PurchasePrice: item.PurchasePrice,
			CategoryName:  item.Category.Name,
			UnitName:      item.Unit.Name,
			UnitSymbol:    item.Unit.Symbol,
		})
	}
	return results, nil
}

func FetchInactiveItemResources(db *gorm.DB, shopId string) ([]ItemResource, error) {
	var items []Item
	if err := db.Preload("Category").Preload("Unit").Where("shop_id = ? AND is_active = 0", shopId).Find(&items).Error; err != nil {
		return nil, err
	}
	results := make([]ItemResource, 0, len(items))
	for _, item := range items {
		results = append(results, ItemResource{
			Id:            item.Id,
			Name:          item.Name,
			SalesPrice:    item.SalesPrice,
			PurchasePrice: item.PurchasePrice,
			CategoryName:  item.Category.Name,
			UnitName:      item.Unit.Name,
			UnitSymbol:    item.Unit.Symbol,
		})
	}
	return results, nil
}

func (readers *ReadServices) AfterUpdateItem(db *gorm.DB, shopId string, id int) error {
	if err := readers.ItemGetService.CleanCache(id); err != nil {
		return err
	}
	if err := readers.ItemListService.CleanCache(shopId); err != nil {
		return err
	}
	return nil
}

type ItemSearch struct {
	Search           string
	CategoryId       int
	UnitId           int
	SalesPriceMin    *decimal.Decimal
	SalesPriceMax    *decimal.Decimal
	PurchasePriceMin *decimal.Decimal
	PurchasePriceMax *decimal.Decimal
}

func SearchItems(db *gorm.DB, shopId string, search *ItemSearch) ([]ItemResource, error) {
	var items []Item
	dbCtx := db.Where("shop_id = ?", shopId)
	if search.Search != "" {
		dbCtx.Where("name LIKE ? OR description LIKE ?",
			"%"+search.Search+"%",
			"%"+search.Search+"%",
		)
	}
	if search.CategoryId > 0 {
		dbCtx.Where("category_id = ?", search.CategoryId)
	}
	if search.UnitId > 0 {
		dbCtx.Where("unit_id = ?", search.UnitId)
	}
	if search.SalesPriceMin != nil {
		if search.SalesPriceMax == nil {
			dbCtx.Where("sales_price >= ?", search.SalesPriceMin)
		} else {
			dbCtx.Where("sales_price >= ? AND sales_price <= ?", search.SalesPriceMin, search.SalesPriceMax)
		}
	}
	if search.PurchasePriceMin != nil {
		if search.PurchasePriceMax == nil {
			dbCtx.Where("purchase_price >= ?", search.PurchasePriceMin)
		} else {
			dbCtx.Where("purchase_price >= ? AND purchase_price <= ?", search.PurchasePriceMin, search.PurchasePriceMax)
		}
	}

	if err := dbCtx.Preload("Category").Preload("Unit").Find(&items).Error; err != nil {
		return nil, err
	}

	results := make([]ItemResource, 0, len(items))
	for _, item := range items {
		results = append(results, ItemResource{
			Id:            item.Id,
			Name:          item.Name,
			SalesPrice:    item.SalesPrice,
			PurchasePrice: item.PurchasePrice,
			CategoryName:  item.Category.Name,
			UnitName:      item.Unit.Name,
			UnitSymbol:    item.Unit.Symbol,
		})
	}
	return results, nil
}
