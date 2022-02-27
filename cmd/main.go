package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/omegabytes/services-api/models"
)

type handler struct {
	db *sql.DB
}

var config struct {
	Port           uint16
	RequestTimeout uint16
	psqlUsername   string
	psqlPassword   string
	psqlHost       string
	psqlPort       uint16
	psqlDatabase   string
}

func init() {
	config.Port = 8080
	config.RequestTimeout = 10
	config.psqlUsername = "postgres"
	config.psqlPassword = "pass"
	config.psqlHost = "services-psql"
	config.psqlPort = 5432
	config.psqlDatabase = "services"
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
	limit := 12                  // todo: user-defined limits to support dynamic page sizes in the ui

	queryStmt := fmt.Sprintf("SELECT * FROM servicetable LIMIT %d", limit)

	if offset != "" {
		queryStmt = fmt.Sprintf("%s OFFSET %s", queryStmt, offset)
	}

	rows, err := h.db.Query(queryStmt)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	results := scanResults(rows)
	encode, _ := json.Marshal(results)
	w.Write(encode)
}

func (h *handler) SearchServiceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call.SearchService")
	vars := r.URL.Query()
	searchTerm := vars.Get("search")
	fmt.Println(searchTerm)

	queryStmt := fmt.Sprintf(`SELECT * FROM servicetable WHERE SIMILARITY((name || ' ' || description), '%s') > 0.12;`, searchTerm)
	rows, err := h.db.Query(queryStmt)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	results := scanResults(rows)
	encode, _ := json.Marshal(results)
	w.Write(encode)
}

func (h *handler) GetServiceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call.GetService")
	vars := mux.Vars(r)

	requestedId, ok := vars["id"]
	if !ok {
		panic("id must not be nil")
	}

	results := []models.Service{}
	// todo: sqlinjection guard
	// tested on curl localhost:8080/services/105%20or%201%3D1
	row := h.db.QueryRow("SELECT * FROM servicetable WHERE id = $1;", requestedId)

	var id int
	var description sql.NullString
	var name string

	switch err := row.Scan(&id, &name, &description); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		s := models.Service{
			Id:          id,
			Name:        name,
			Description: description.String,
		}
		results = append(results, s)
	default:
		panic(err)
	}

	encode, _ := json.Marshal(results)

	w.Write(encode)
}

func scanResults(rows *sql.Rows) []models.Service {
	results := []models.Service{}
	for rows.Next() {
		var id int
		var description sql.NullString
		var name string

		switch err := rows.Scan(&id, &name, &description); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			s := models.Service{
				Id:          id,
				Name:        name,
				Description: description.String,
			}
			results = append(results, s)
		default:
			panic(err)
		}
	}

	err := rows.Err()
	if err != nil {
		panic(err)
	}
	return results
}

func connectToDB(uri string) (*sql.DB, error) {
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10)

	return db, nil
}

func buildConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.psqlHost, config.psqlPort, config.psqlUsername, config.psqlPassword, config.psqlDatabase)
}
