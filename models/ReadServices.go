package models

import (
	"linn221/shop/services"
	"time"

	"gorm.io/gorm"
)

type ReadServices struct {
	CategoryGetService  services.Getter[CategoryDetailResource]
	CategoryListService services.Lister[CategoryResource]
	UnitGetService      services.Getter[UnitDetailResource]
	UnitListService     services.Lister[UnitResource]
	ItemGetService      services.Getter[ItemDetailResource]
	ItemListService     services.Lister[ItemResource]
}

func NewReaders(db *gorm.DB, cache services.CacheService) *ReadServices {
	return &ReadServices{
		CategoryGetService: &defaultGetService[CategoryDetailResource]{
			db:          db,
			cache:       cache,
			table:       "categories",
			cachePrefix: "Category",
			cacheLength: forever,
		},
		CategoryListService: &defaultListService[CategoryResource]{
			db:          db,
			cache:       cache,
			table:       "categories",
			cachePrefix: "CategoryList",
			cacheLength: forever,
		},
		UnitGetService: &defaultGetService[UnitDetailResource]{
			db:          db,
			cache:       cache,
			table:       "units",
			cachePrefix: "Unit",
			cacheLength: forever,
		},
		UnitListService: &defaultListService[UnitResource]{
			db:          db,
			cache:       cache,
			table:       "units",
			cachePrefix: "UnitList",
			cacheLength: forever,
		},
		ItemGetService: &customGetService[ItemDetailResource]{
			db:          db,
			cache:       cache,
			cachePrefix: "items",
			cacheLength: time.Hour * 127,
			fetch: func(db *gorm.DB, id int) (ItemDetailResource, error) {
				var item Item
				if err := db.Preload("Category").Preload("Unit").First(&item, id).Error; err != nil {
					return ItemDetailResource{}, err
				}
				result := ItemDetailResource{
					Id:            item.Id,
					Name:          item.Name,
					SalesPrice:    item.SalesPrice,
					PurchasePrice: item.PurchasePrice,
					Description:   item.Description,
				}
				result.ShopId = item.ShopId
				result.IsActive = item.IsActive
				result.Category.Id = item.Category.Id
				result.Category.Name = item.Category.Name
				result.Category.IsActive = item.Category.IsActive

				result.Unit.Id = item.Unit.Id
				result.Unit.Name = item.Unit.Name
				result.Unit.Symbol = item.Unit.Symbol
				result.Unit.IsActive = item.Unit.IsActive
				return result, nil
			},
		},
		ItemListService: &customListService[ItemResource]{
			db:          db,
			cache:       cache,
			cachePrefix: "ItemList",
			cacheLength: forever,
			fetch:       FetchItemResources,
		},
	}
}
