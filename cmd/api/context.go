package api

import (
	"context"
	"net/http"
	"rest-api/internal/resource/model"
)

type contextKey string

const userContextKey = contextKey("user")

func ContextSetUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func ContextGetUser(r *http.Request) *model.User {
	user, ok := r.Context().Value(userContextKey).(*model.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
