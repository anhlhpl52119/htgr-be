package app

import (
	"database/sql"
	"htrr-apis/internal/api"
	"htrr-apis/internal/store"
	"htrr-apis/internal/utils"
	"htrr-apis/migrations"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger            *log.Logger
	DB                *sql.DB
	UserHandler       *api.UserHandler
	RestaurantHandler *api.RestaurantHandler
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFs(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "[APP] ", log.Ldate|log.Ltime)

	userHandler := api.NewUserHandler(
		logger, store.NewPostgresUserStore(pgDB))

	restaurantHandler := api.NewRestaurantHandler(
		logger,
		store.NewPostgresRestaurantStore(pgDB))

	app := &Application{
		Logger:            logger,
		DB:                pgDB,
		UserHandler:       userHandler,
		RestaurantHandler: restaurantHandler,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Available!!!"})
}
