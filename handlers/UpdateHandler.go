package handlers

import (
	"context"
	"linn221/shop/myctx"
	"linn221/shop/views"
	"net/http"
	"strconv"
)

type UpdateResource[T any] struct {
	parseInput       func(r *http.Request) (*T, formErrors)
	handleParseError func(ctx context.Context, r *http.Request, session *Session, input *T, fe formErrors, renderer *views.Renderer) error
	handle           func(ctx context.Context, r *http.Request, session *Session, input *T, renderer *views.Renderer) error
}

func UpdateHandler[T any](t *views.Templates, res UpdateResource[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, err := myctx.GetUserId(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resIdStr := r.PathValue("id")
		if resIdStr == "" {
			http.Error(w, "resource id is required", http.StatusBadRequest)
			return
		}
		resId, err := strconv.Atoi(resIdStr)
		if err != nil {
			http.Error(w, "resource id is required", http.StatusBadRequest)
			return

		}
		session := Session{
			UserId: userId,
			ResId:  resId,
		}

		renderer := t.NewRenderer(w, userId)
		input, ferrs := res.parseInput(r)
		if len(ferrs) > 0 {
			w.Header().Add("HX-Retarget", "#edit-form")
			w.Header().Add("HX-Reswap", "outerHTML")
			finalErrHandle(w,
				res.handleParseError(r.Context(), r, &session, input, ferrs, renderer),
			)
			return
		}

		err = res.handle(ctx, r, &session, input, renderer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// type UpdateFunc[T any] func(w http.ResponseWriter, r *http.Request, session Session, input *T) error

// func UpdateHandler[T any](handle UpdateFunc[T]) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		ctx := r.Context()
// 		userId, shopId, err := myctx.GetIdsFromContext(ctx)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		input, ok, err := parseJson[T](w, r)
// 		if !ok {
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 			}
// 			return
// 		}

// 		UpdateSession := Session{
// 			UserId: userId,
// 		}
// 		err = handle(w, r, UpdateSession, input)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

// type UpdateFunc[T any] func(w http.ResponseWriter, r *http.Request, session Session, input *T) error

// func UpdateHandler[T any](handle UpdateFunc[T]) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		ctx := r.Context()
// 		userId, shopId, err := myctx.GetIdsFromContext(ctx)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		input, ok, err := parseJson[T](w, r)
// 		if !ok {
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 			}
// 			return
// 		}

// 		resIdStr := r.PathValue("id")
// 		resId, err := strconv.Atoi(resIdStr)
// 		if err != nil || resId <= 0 {
// 			http.Error(w, "invalid resource id", http.StatusBadRequest)
// 			return
// 		}

// 		UpdateSession := Session{
// 			UserId: userId,
// 			ShopId: shopId,
// 			ResId:  resId,
// 		}
// 		err = handle(w, r, UpdateSession, input)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }
