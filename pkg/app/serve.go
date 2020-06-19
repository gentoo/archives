// Entrypoint for the web application

package app

import (
	"archives/pkg/app/home"
	"archives/pkg/app/list"
	"archives/pkg/app/message"
	"archives/pkg/app/popular"
	"archives/pkg/app/search"
	"archives/pkg/config"
	"fmt"
	"log"
	"net/http"
)

// Serve is used to serve the web application
func Serve() {

	fmt.Println("Serving on Port " + config.Port())

	for _, mailingList := range config.AllPublicMailingLists() {
		setRoute("/"+mailingList+"/message/", message.Show)
		setRoute("/"+mailingList+"/messages/", list.Messages)
		setRoute("/"+mailingList+"/threads/", list.Threads)
		setRoute("/"+mailingList+"/", list.Show)
	}

	setRoute("/lists", list.Browse)

	setRoute("/popular", popular.Threads)

	setRoute("/search", search.Search)

	setRoute("/", home.Show)

	fs := http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))
	http.Handle("/assets/", fs)

	log.Fatal(http.ListenAndServe(":"+config.Port(), nil))

}

// define a route using the default middleware and the given handler
func setRoute(path string, handler http.HandlerFunc) {
	http.HandleFunc(path, mw(handler))
}

// mw is used as default middleware to set the default headers
func mw(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setDefaultHeaders(w)
		handler(w, r)
	}
}

// setDefaultHeaders sets the default headers that apply for all pages
func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", config.CacheControl())
}
