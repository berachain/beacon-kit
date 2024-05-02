package main

import (
	"fmt"
	"net/http"
	apihandlers "github.com/berachain/beacon-kit/mod/api/handlers"
	chi "github.com/go-chi/chi/v5"
	cors "github.com/go-chi/cors"
)

func main() {

	arg := "chi"
	var r apihandlers.Router
	switch arg {
	case "chi":
		r = chi.NewRouter()
	}

	corsMiddleware := cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	apihandlers.UseMiddlewares(r, []func(next http.Handler) http.Handler{corsMiddleware})
	apihandlers.AssignRoutes(r, apihandlers.RouteHandler{})
	fmt.Println("Server starting on port 3000")
	http.ListenAndServe(":3000", r)
}
