package handlers

import (
	"linn221/shop/myctx"
	"net/http"
	"strconv"
)

type UpdateFunc[T any] func(w http.ResponseWriter, r *http.Request, session Session, input *T) error

func UpdateHandler[T any](handle UpdateFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userId, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		input, ok, err := parseJson[T](w, r)
		if !ok {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		resIdStr := r.PathValue("id")
		resId, err := strconv.Atoi(resIdStr)
		if err != nil || resId <= 0 {
			http.Error(w, "invalid resource id", http.StatusBadRequest)
			return
		}

		UpdateSession := Session{
			UserId: userId,
			ShopId: shopId,
			ResId:  resId,
		}
		err = handle(w, r, UpdateSession, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
