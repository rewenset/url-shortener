package main

import (
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

var last int
var urls = make(map[int]string)

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
		}
		last++
		urls[last] = r.PostFormValue("original")
		data = strconv.Itoa(last)
	}

	if err := t.Execute(w, data); err != nil {
		log.Printf("could not execute template: %v", err)
	}

}

func follow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	urlID, err := strconv.Atoi(vars["urlID"])
	if err != nil {
		log.Printf("could not parse urlID to int: %v", err)
		http.Error(w, "bad url", http.StatusBadRequest)
		return
	}

	if e, ok := urls[urlID]; ok {
		log.Printf("redirecting to %s", e)
		http.Redirect(w, r, e, http.StatusFound)
	} else {
		log.Printf("requested urlID - %d - not found", urlID)
		http.NotFound(w, r)
	}
}
