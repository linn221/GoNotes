package handlers

import (
	"linn221/shop/myctx"
	"linn221/shop/services"
	"net/http"
	"strconv"
)

type GetSession struct {
	UserId int
	ShopId string
	ResId  int
}

type GetFunc func(w http.ResponseWriter, r *http.Request, session *GetSession) error

func DefaultGetHandler[T any](getService services.Getter[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		_, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resIdStr := r.PathValue("id")
		resId, err := strconv.Atoi(resIdStr)
		if err != nil || resId <= 0 {
			http.Error(w, "invalid resource id", http.StatusBadRequest)
			return
		}

		resource, found, err := getService.Get(shopId, resId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !found {
			finalErrHandle(w,
				respondNotFound(w, "record not found"),
			)
			return
		}

		finalErrHandle(w,
			respondOk(w, resource),
		)
	}
}

func finalErrHandle(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
