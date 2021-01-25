package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.GET("/", index)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "index.html")
}