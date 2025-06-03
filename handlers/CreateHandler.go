package handlers

import (
	"linn221/shop/myctx"
	"net/http"
)

type CreateFunc[T any] func(w http.ResponseWriter, r *http.Request, session Session, input *T) error

func CreateHandler[T any](handle CreateFunc[T]) http.HandlerFunc {
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

		CreateSession := Session{
			UserId: userId,
			ShopId: shopId,
		}
		err = handle(w, r, CreateSession, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
