package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/omegabytes/services-api/models"
)

type handler struct {
	db *sql.DB
}

var config struct {
	Port           uint16
	RequestTimeout uint16
	MySQLUsername  string
	MySQLPassword  string
	MySQLHost      string
	MySQLPort      uint16
	MySQLDatabase  string
}

func init() {
	config.Port = 8080
	config.RequestTimeout = 10
	config.MySQLUsername = "user"
	config.MySQLPassword = "pass"
	config.MySQLHost = "services-mysql"
	config.MySQLPort = 3306
	config.MySQLDatabase = "services"
}

func main() {
	fmt.Println(fmt.Sprintf("server started at port %d", config.Port))

	fmt.Println(fmt.Sprintf("connecting to %s on port %d", config.MySQLDatabase, config.MySQLPort))
	db, err := connectToDB(buildConnectionString())
	if err != nil {
		log.Fatal("db err", err)
	}
	defer db.Close()
	fmt.Println(fmt.Sprintf("connected to %s", config.MySQLDatabase))

	h := handler{db}
	r := mux.NewRouter()
	r.HandleFunc("/services/{id}", h.GetServiceHandler)
	r.HandleFunc("/services", h.ListServiceHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}

func (h *handler) ListServiceHandler(w http.ResponseWriter, r *http.Request) {
	sqlStatement := `SELECT * FROM servicetable;`
	results := []models.Service{}

	rows, err := h.db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var description sql.NullString
		var name string

		switch err = rows.Scan(&id, &name, &description); err {
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

	err = rows.Err()
	if err != nil {
		panic(err)
	}
	encode, _ := json.Marshal(results)

	w.Write(encode)
}

func (h *handler) GetServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	requestedId, ok := vars["id"]
	if !ok {
		panic("id must not be nil")
	}

	results := []models.Service{}
	// todo: sqlinjection guard
	// tested on curl localhost:8080/services/105%20or%201%3D1
	row := h.db.QueryRow("SELECT * FROM servicetable WHERE id = ?;", requestedId)

	var id string
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

func connectToDB(uri string) (*sql.DB, error) {
	db, err := sql.Open("mysql", uri)
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
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.MySQLUsername, config.MySQLPassword, config.MySQLHost, config.MySQLPort, config.MySQLDatabase)
}
