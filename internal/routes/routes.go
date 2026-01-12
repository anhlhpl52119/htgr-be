package routes

import (
	"htrr-apis/internal/app"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// health
	r.Get("/health", app.HealthCheck)

	// user
	r.Post("/user", app.UserHandler.HandleCreateUser)

	// restaurants
	r.Get("/restaurants", app.RestaurantHandler.HandleSearchRestaurant)
	r.Post("/restaurant", app.RestaurantHandler.HandleCreateRestaurant)
	r.Get("/restaurant/{id}", app.RestaurantHandler.HandleGetRestaurantById)
	r.Patch("/restaurant/{id}", app.RestaurantHandler.HandleUpdateRestaurant)
	r.Delete("/restaurant/{id}", app.RestaurantHandler.HandleDeleteRestaurant)
	r.Delete("/restaurants", app.RestaurantHandler.HandleBulkDeleteRestaurants)
	return r
}
