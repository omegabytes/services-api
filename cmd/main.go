package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/omegabytes/services-api/handlers"
	"github.com/omegabytes/services-api/store"
)

var config struct {
	// service
	Port           uint16
	limit          uint16
	precision      float32
	RequestTimeout uint16

	// database
	psqlUsername    string
	psqlPassword    string
	psqlHost        string
	psqlPort        uint16
	psqlDatabase    string
	maxIdleConns    int
	maxOpenConns    int
	maxConnLifetime time.Duration
}

func init() {
	config.Port = 8080
	config.limit = 12       // todo: user-defined limits to support dynamic page sizes in the ui
	config.precision = 0.12 // magic number derived from observing search behavior on my generated dataset
	config.RequestTimeout = 10

	config.psqlUsername = "postgres"
	config.psqlPassword = "pass"
	config.psqlHost = "services-psql"
	config.psqlPort = 5432
	config.psqlDatabase = "services"
	config.maxIdleConns = 1
	config.maxOpenConns = 10
	config.maxConnLifetime = 10
}

func main() {
	fmt.Println(fmt.Sprintf("server started at port %d", config.Port))

	db, err := connectToDB(buildConnectionString())
	if err != nil {
		log.Fatal("db err", err)
	}
	defer db.Close()
	fmt.Println(fmt.Sprintf("connected to %s", config.psqlDatabase))

	h := handlers.Handler{
		Store: store.Store{
			DB:        db,
			Limit:     config.limit,
			Precision: config.precision,
		},
	}

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Use(h.Middleware)
	api := r.PathPrefix("/api").Subrouter() // some /api paths do not require authentication
	api.HandleFunc("/authenticate", h.Authenticate).Methods("POST")

	v1 := api.PathPrefix("/v1").Subrouter()
	v1.Use(h.AuthMiddleware) // all requests to /api/v1 paths require authentication
	v1.HandleFunc("/services", h.SearchServiceHandler).Queries("search", "{search}").Methods("GET")
	v1.HandleFunc("/services", h.ListServiceHandler).Methods("GET")
	v1.HandleFunc("/services/{id}", h.GetServiceHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}

func connectToDB(uri string) (*sql.DB, error) {
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(config.maxIdleConns)
	db.SetMaxOpenConns(config.maxOpenConns)
	db.SetConnMaxLifetime(config.maxConnLifetime * time.Second)

	return db, nil
}

func buildConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.psqlHost, config.psqlPort, config.psqlUsername, config.psqlPassword, config.psqlDatabase)
}
