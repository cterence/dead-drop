package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/cterence/dead-drop/views"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	// Check if env is development
	indexComponent := views.Index()
	dropComponent := views.GetDrop()

	dbHost, isSet := os.LookupEnv("DB_HOST")
	if !isSet {
		dbHost = "localhost"
	}
	dbPort, isSet := os.LookupEnv("DB_PORT")
	if !isSet {
		dbPort = "8080"
	}

	dbUrl := "http://" + dbHost + ":" + dbPort

	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create table drops
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS drops (id TEXT PRIMARY KEY, data TEXT)")
	if err != nil {
		panic(err)
	}

	http.Handle("/", templ.Handler(indexComponent))

	http.Handle("/drop/put", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Receive the data.
		r.ParseForm()
		data := r.Form.Get("data")
		// Create random link with the data.

		randomBytes := make([]byte, 32)
		_, err := rand.Read(randomBytes)
		if err != nil {
			panic(err)
		}
		randomId := base64.URLEncoding.EncodeToString(randomBytes)[:12]

		// protocol is https if the proxy header is set, otherwise http.
		protocol := "http://"
		if r.Header.Get("X-Forwarded-Proto") == "https" {
			protocol = "https://"
		}

		host := r.Host

		link := protocol + host + "/get/" + randomId
		fmt.Println("Link:", link)
		fmt.Println("Data:", data)

		_, err = db.Exec("INSERT INTO drops (id, data) VALUES (?, ?)", randomId, data)
		if err != nil {
			panic(err)
		}

		// Send the link back to the user.
		w.Write([]byte(link))
	}))

	http.Handle("/get/", templ.Handler(dropComponent))

	http.Handle("/drop/get/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the ID from the URL.
		id := r.URL.Path[len("/drop/get/"):]

		// Get the data from the database.
		rows, err := db.Query("SELECT data FROM drops WHERE id = ?", id)
		if err != nil {
			panic(err)
		}

		var data string
		for rows.Next() {
			err = rows.Scan(&data)
			if err != nil {
				panic(err)
			}
		}

		_, err = db.Exec("DELETE FROM drops WHERE id = ?", id)
		if err != nil {
			panic(err)
		}

		w.Write([]byte(data))
	}))

	fmt.Println("Listening on :3000")
	http.ListenAndServe("0.0.0.0:3000", nil)
}
