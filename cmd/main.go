package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/omegabytes/services-api/models"
)

type handler struct {
	db *sql.DB
}

var config struct {
	// service
	Port      uint16
	limit     uint16
	precision float32

	// database
	RequestTimeout  uint16
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

	fmt.Println(fmt.Sprintf("connecting to %s on port %d", config.psqlDatabase, config.psqlPort))
	db, err := connectToDB(buildConnectionString())
	if err != nil {
		log.Fatal("db err", err)
	}
	defer db.Close()
	fmt.Println(fmt.Sprintf("connected to %s", config.psqlDatabase))

	h := handler{db}
	r := mux.NewRouter()
	r.HandleFunc("/services", h.SearchServiceHandler).Queries("search", "{search}")
	r.HandleFunc("/services", h.ListServiceHandler)
	r.HandleFunc("/services/{id}", h.GetServiceHandler)
	http.Handle("/", r)

	// go func() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
	// }()

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// <-c
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RequestTimeout*uint16(time.Second)))
	// defer cancel()
	// srv.Shutdown(ctx)
	// log.Println("shutting down")
	// os.Exit(0)

}

func (h *handler) ListServiceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call.ListService")
	vars := r.URL.Query()
	offset := vars.Get("offset") // assume offset = last record shown + 1, handled by the front end

	queryStmt := fmt.Sprintf("SELECT * FROM servicetable LIMIT %d", config.limit)

	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			http.Error(w, "Invalid offset", 400)
			return
		}
		queryStmt = fmt.Sprintf("%s OFFSET %d", queryStmt, o)
	}

	rows, err := h.db.Query(queryStmt)
	if err != nil {
		http.Error(w, "Invalid query", 500)
	}
	defer rows.Close()

	results, err := scanResults(rows)
	if err != nil {
		http.Error(w, "Error scanning the database", 500)
	}
	encode, _ := json.Marshal(results)
	w.Write(encode)
}

func (h *handler) SearchServiceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call.SearchService")
	vars := r.URL.Query()
	searchTerm := vars.Get("search")

	queryStmt := fmt.Sprintf(`SELECT * FROM servicetable WHERE SIMILARITY((name || ' ' || description), '%s') > %f limit %d;`, searchTerm, config.precision, config.limit)
	rows, err := h.db.Query(queryStmt)
	if err != nil {
		http.Error(w, "Invalid query", 500)
	}
	defer rows.Close()

	results, err := scanResults(rows)
	if err != nil {
		http.Error(w, "Error scanning the database", 500)
		return
	}

	encode, _ := json.Marshal(results)
	w.Write(encode)
}

func (h *handler) GetServiceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call.GetService")
	vars := mux.Vars(r)

	requestedId, ok := vars["id"]
	if !ok {
		http.Error(w, "'id' is required", 400)
	}

	results := []models.Service{}
	// todo: sqlinjection guard
	// tested on curl localhost:8080/services/105%20or%201%3D1
	row := h.db.QueryRow("SELECT * FROM servicetable WHERE id = $1;", requestedId)

	var id int
	var description sql.NullString
	var name string

	switch err := row.Scan(&id, &name, &description); err {
	case nil:
		s := models.Service{
			Id:          id,
			Name:        name,
			Description: description.String,
		}
		results = append(results, s)
	default:
		http.Error(w, "Error scanning the database", 500)
	}

	encode, _ := json.Marshal(results)
	w.Write(encode)
}

func scanResults(rows *sql.Rows) ([]models.Service, error) {
	results := []models.Service{}
	for rows.Next() {
		var id int
		var description sql.NullString
		var name string

		switch err := rows.Scan(&id, &name, &description); err {
		case nil:
			s := models.Service{
				Id:          id,
				Name:        name,
				Description: description.String,
			}
			results = append(results, s)
		default:
			return nil, err
		}
	}

	err := rows.Err()
	if err != nil {
		return nil, err
	}
	return results, nil
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
