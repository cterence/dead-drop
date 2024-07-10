package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/cterence/dead-drop/views"
)

func main() {
	component := views.Index()

	http.Handle("/", templ.Handler(component))

	http.Handle("/data", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Receive the data.
		r.ParseForm()
		// data := r.Form.Get("data")
	}))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
