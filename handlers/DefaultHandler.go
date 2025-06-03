package handlers

import (
	"linn221/shop/myctx"
	"linn221/shop/services"
	"net/http"
)

type DefaultSession struct {
	UserId int
	ShopId string
}

func DefaultHandler(handle func(http.ResponseWriter, *http.Request, *DefaultSession) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session := DefaultSession{
			UserId: userId,
			ShopId: shopId,
		}

		err = handle(w, r, &session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func DefaultListHandler[T any](listService services.Lister[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resourceSlice, err := listService.List(shopId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		finalErrHandle(w,
			respondOk(w, resourceSlice),
		)
	}
}
