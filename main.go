package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/Blustak/go-transActor/internal/config"
	"github.com/Blustak/go-transActor/internal/database"
	"github.com/Blustak/go-transActor/internal/repl"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	godotenv.Load()
	dbPath := handleEnvLookup("DB_PATH")
	sqlDriver := handleEnvLookup("SQL_DRIVER")
	dbConn, err := sql.Open(sqlDriver, dbPath)
	if err != nil {
		panic(err)
	}
	r := os.Stdin
	w := os.Stdout
	progCfg := config.Config{
		Log:    slog.New(slog.NewTextHandler(w, nil)),
		DB:     database.New(dbConn),
		Reader: r,
		Writer: w,
	}
	repler := repl.NewRepler(&progCfg)
	repler.Run()
}

func handleEnvLookup(k string) string {
	v, ok := os.LookupEnv(k)
	if !ok {
		panic(fmt.Sprintf("Couldn't find env variable %s\n", k))
	}
	return v
}
