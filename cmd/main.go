package main

import (
	"fmt"
	"log"
	"net/http"
)

var config struct {
	Port           uint16
	RequestTimeout uint16
}

func init() {
	config.Port = 8080
	config.RequestTimeout = 10
}

func main() {

	http.HandleFunc("/", HelloHandler)

	fmt.Println(fmt.Sprintf("server started at port %d", config.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := "guest"
	keys, ok := r.URL.Query()["name"]
	if ok {
		name = keys[0]
	}
	fmt.Fprintf(w, "hello %s\n", name)
}

// get services
