package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var indexTmpl = template.Must(template.ParseFiles("index.html"))

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

	if err := indexTmpl.Execute(w, data); err != nil {
		log.Printf("could not execute template: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

}

func follow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if url := load(vars["urlID"]); url != "" {
		log.Printf("redirecting to %s", url)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		log.Printf("requested urlID - %s - not found", vars["urlID"])
		http.NotFound(w, r)
	}
}
