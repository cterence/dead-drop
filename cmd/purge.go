/*
Copyright © 2024 Térence Chateigné
*/
package cmd

import (
	"database/sql"
	"log/slog"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbHost := viper.GetString("db_host")
		dbPort := viper.GetString("db_port")
		dbUrl := "http://" + dbHost + ":" + dbPort

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

		// Delete all drops with a timestamp older than 24 hours
		res, err := db.Exec("DELETE FROM drops WHERE timestamp < datetime('now', '-1 day')")
		if err != nil {
			slog.Error("Failed to purge drops")
			panic(err)
		}
		purgedDropsInt, err := res.RowsAffected()
		if err != nil {
			slog.Error("Failed to get the number of purged drops")
			panic(err)
		}
		purgedDrops := strconv.Itoa(int(purgedDropsInt))

		slog.Info("Drops purged: " + purgedDrops)
	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)

	flags := purgeCmd.Flags()
	flags.SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ReplaceAll(name, "-", "_"))
	})

	flags.String("db-host", "localhost", "The database host")
	flags.String("db-port", "8080", "The database port")

	viper.BindPFlags(flags)
}
