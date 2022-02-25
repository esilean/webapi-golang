package main

import (
	"backend/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	jwt struct {
		secret string
	}
	theMovieDB struct {
		key string
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func init() {
	godotenv.Load()
}

func main() {
	var cfg config

	cfg.port, _ = strconv.Atoi(os.Getenv("GO_MOVIES_PORT"))
	cfg.env = os.Getenv("GO_MOVIES_ENV")
	cfg.db.dsn = os.Getenv("GO_MOVIES_POSTGRES_DSN")
	cfg.jwt.secret = os.Getenv("GO_MOVIES_JWT_SECRET")
	cfg.theMovieDB.key = os.Getenv("GO_MOVIES_THEMOVIEDB_KEY")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	app := application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", cfg.port)
	err = srv.ListenAndServe()
	if err != nil {
		app.logger.Println(err)
	}
}

func openDB(cfg config) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(cfg.db.dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
