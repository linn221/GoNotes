package models

import (
	"linn221/shop/services"

	"gorm.io/gorm"
)

type ReadServices struct {
	// CategoryGetService  services.Getter[CategoryDetailResource]
	// CategoryListService services.Lister[CategoryResource]
	// UnitGetService      services.Getter[UnitDetailResource]
	// UnitListService     services.Lister[UnitResource]
	// ItemGetService      services.Getter[ItemDetailResource]
	// ItemListService     services.Lister[ItemResource]
	NoteService NoteService
}

func NewReaders(db *gorm.DB, cache services.CacheService) *ReadServices {
	return &ReadServices{
		NoteService: NoteService{db: db},
	}
}
