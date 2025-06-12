package models

import (
	"linn221/shop/services"

	"gorm.io/gorm"
)

type CrudServices struct {
	// CategoryGetService  services.Getter[CategoryDetailResource]
	// CategoryListService services.Lister[CategoryResource]
	// UnitGetService      services.Getter[UnitDetailResource]
	// UnitListService     services.Lister[UnitResource]
	// ItemGetService      services.Getter[ItemDetailResource]
	// ItemListService     services.Lister[ItemResource]
	NoteService  *NoteService
	LabelService *LabelService
}

func NewServices(db *gorm.DB, cache services.CacheService) *CrudServices {
	return &CrudServices{
		NoteService:  &NoteService{db: db},
		LabelService: &LabelService{db: db},
	}
}
