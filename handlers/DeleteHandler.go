package handlers

import (
	"linn221/shop/myctx"
	"net/http"
	"strconv"
)

type Session struct {
	UserId int
	ShopId string
	ResId  int
}

type DeleteFunc func(w http.ResponseWriter, r *http.Request, session Session) error

func DeleteHandler(handle DeleteFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userId, shopId, err := myctx.GetIdsFromContext(ctx)
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

		DeleteSession := Session{
			UserId: userId,
			ShopId: shopId,
			ResId:  resId,
		}
		err = handle(w, r, DeleteSession)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
