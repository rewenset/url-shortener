package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", logger(index))
	r.HandleFunc("/f/{urlID}", logger(follow))

	log.Println("starting URL shortener server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func logger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		h(w, r)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("index.html"))

	var data string
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Printf("could not parse form: %v", err)
		} else {
			urlID := save(r.PostFormValue("original"))
			data = fmt.Sprintf("http://%s/f/%s", r.Host, urlID)
		}
	}

	if err := t.Execute(w, data); err != nil {
		log.Printf("could not execute template: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

}

func follow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if original := load(vars["urlID"]); original != "" {
		log.Printf("redirecting to %s", original)
		http.Redirect(w, r, original, http.StatusFound)
	} else {
		log.Printf("requested urlID - %s - not found", vars["urlID"])
		http.NotFound(w, r)
	}
}
