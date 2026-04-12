package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/cors"
)

func CORS() func(http.Handler) http.Handler {
	origins := []string{
		"https://danielkrammes.com",
		"https://www.danielkrammes.com",
	}
	// Dev: localhost ports 3000-3010
	for port := 3000; port <= 3010; port++ {
		origins = append(origins, fmt.Sprintf("http://localhost:%d", port))
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Tenant-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
