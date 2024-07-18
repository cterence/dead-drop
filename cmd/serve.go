/*
Copyright © 2024 Térence Chateigné

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/cterence/dead-drop/views"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Create a dead-drop server",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbHost := viper.GetString("db_host")
		dbPort := viper.GetString("db_port")
		address := viper.GetString("address")
		port := viper.GetString("port")

		dbUrl := "http://" + dbHost + ":" + dbPort

		// Check if env is development
		indexComponent := views.Index()
		dropComponent := views.GetDrop()

		db, err := sql.Open("libsql", dbUrl)
		if err != nil {
			slog.Error("Failed to open database connection")
			panic(err)
		}
		defer db.Close()

		// Check if the database is up.
		_, err = db.Exec("SELECT 1")
		if err != nil {
			slog.Error("Cannot connect to the database")
			return
		}

		// Check if the table drops exists
		_, err = db.Exec("SELECT * FROM drops LIMIT 1")
		if err != nil {
			slog.Error("Table drops does not exist")
			panic(err)
		}

		http.Handle("/", templ.Handler(indexComponent))

		http.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the database is up.
			// db.Ping() does not work for some reason.
			http.Header.Add(w.Header(), "Content-Type", "application/json")
			_, err := db.Exec("SELECT 1")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status": {"server": "OK", "database": "KO"}}`))
				slog.Error("Database is down")
			} else {
				slog.Info("Database is up")
				// Write OK to the response as a JSON object.
				w.Write([]byte(`{"status": {"server": "OK", "database": "OK"}}`))
			}
		}))

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

			// Save the data to the database.
			// Database : id, timestamp, data
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

		slog.Info("Listening on " + address + ":" + port)
		err = http.ListenAndServe(address+":"+port, nil)
		if err != nil {
			slog.Error("Failed to start server")
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	flags := serveCmd.Flags()
	flags.SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ReplaceAll(name, "-", "_"))
	})

	flags.String("db-host", "localhost", "The database host")
	flags.String("db-port", "8080", "The database port")
	flags.String("address", "127.0.0.1", "The address to listen on")
	flags.String("port", "3000", "The port to listen on")

	viper.BindPFlags(flags)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
