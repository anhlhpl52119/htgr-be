package app

import (
	"database/sql"
	"fmt"
	"htrr-apis/internal/api"
	"htrr-apis/internal/store"
	"htrr-apis/internal/utils"
	"htrr-apis/migrations"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Application struct {
	Logger            *log.Logger
	DB                *sql.DB
	UserHandler       *api.UserHandler
	RestaurantHandler *api.RestaurantHandler
}

func NewApplication() (*Application, error) {
	// load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (assuming env vars are set)")
	}

	// 1. Try single connection string (Best for Supabase/Prod)
	connStr := os.Getenv("DATABASE_URL")

	// 2. Fallback to individual variables (Keeps Docker setup working)
	if connStr == "" {
		dbName := os.Getenv("DB_NAME")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbHost := os.Getenv("DB_HOST") // Good to have host configurable too
		if dbHost == "" {
			dbHost = "localhost"
		}

		connStr = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbHost, dbUser, dbPassword, dbName, dbPort)
	}

	pgDB, err := store.Open(connStr)
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
