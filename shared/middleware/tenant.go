package middleware

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func RequireTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			http.Error(w, `{"error":"missing X-Tenant-ID header"}`, http.StatusBadRequest)
			return
		}

		if _, err := bson.ObjectIDFromHex(tenantID); err != nil {
			http.Error(w, `{"error":"invalid tenant ID format"}`, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
