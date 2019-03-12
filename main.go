package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", logger(index))
	r.HandleFunc("/f/{urlID}", logger(follow))

	log.Println("Starting URL shortener server on port 8000")
	http.ListenAndServe(":8000", r)
}

func logger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		log.Printf("%s %s", r.Method, r.RequestURI)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "allow creating shortener urls")
}

func follow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlID := vars["urlID"]
	fmt.Fprintf(w, "redirect with 302 or send 404 if unknown for url id: %s", urlID)
}
