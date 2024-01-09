package am

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/adriancrafter/todoapp/internal/am/errors"
)

const (
	TenantIDKey   ContextKey = "tenantID"
	TenantIDParam            = "tenantID"
)

func TenantID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tenantID := vars[TenantIDParam]

		ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetTenantID(r *http.Request) (string, error) {
	tenantID, ok := r.Context().Value(TenantIDKey).(string)
	if !ok {
		return "", errors.NewError("tenantID not found in request")
	}
	return tenantID, nil
}
