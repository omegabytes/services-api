package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

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
	_, err := connectToDB(buildConnectionString())
	if err != nil {
		log.Fatal("db err", err)
	}
	fmt.Println(fmt.Sprintf("connected to %s", config.MySQLDatabase))

	http.HandleFunc("/", HelloHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}

// get services

func GetServiceHandler(w http.ResponseWriter, r *http.Request) {

}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := "guest"
	keys, ok := r.URL.Query()["name"]
	if ok {
		name = keys[0]
	}
	fmt.Fprintf(w, "hello %s\n", name)
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
	defer db.Close()

	return db, nil
}

func buildConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.MySQLUsername, config.MySQLPassword, config.MySQLHost, config.MySQLPort, config.MySQLDatabase)
}
