package handlers

import (
	"context"
	"linn221/shop/myctx"
	"linn221/shop/views"
	"net/http"
)

type CreateResource[T any] struct {
	parseInput       func(*http.Request) (*T, formErrors)
	handleParseError func(ctx context.Context, r *http.Request, session *DefaultSession, input *T, fe formErrors, vr *views.Renderer) error
	handle           func(ctx context.Context, r *http.Request, session *DefaultSession, input *T, vr *views.Renderer) error
}

func CreateHandler[T any](t *views.Templates, res CreateResource[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, err := myctx.GetUserId(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session := DefaultSession{
			UserId: userId,
		}

		renderer := t.NewRenderer(w, userId)
		input, ferrs := res.parseInput(r)
		if len(ferrs) > 0 {
			w.Header().Add("HX-Retarget", "#create-form")
			w.Header().Add("HX-Reswap", "outerHTML")
			finalErrHandle(w,
				res.handleParseError(ctx, r, &session, input, ferrs, renderer),
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

// type CreateFunc[T any] func(w http.ResponseWriter, r *http.Request, session Session, input *T) error

// func CreateHandler[T any](handle CreateFunc[T]) http.HandlerFunc {
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

// 		CreateSession := Session{
// 			UserId: userId,
// 		}
// 		err = handle(w, r, CreateSession, input)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }
