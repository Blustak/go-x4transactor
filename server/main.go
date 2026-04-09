package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Blustak/go-transActor/internal/protocol"
	"github.com/Blustak/go-transActor/server/internal/server"
	"github.com/joho/godotenv"
)

const Hostname = "localhost"
const Port = "8080"

func main() {
	godotenv.Load()
	protocol.Init()
	dbPath, ok := os.LookupEnv("DB_PATH")
	if !ok {
		panic("DB_PATH environment variable not set")
	}
	f := os.Stdout
	if logPath, ok := os.LookupEnv("SERVER_LOG"); ok {
		var err error
		if f, err = os.Open(logPath); err != nil {
			panic(err)
		}
	}
	defer f.Close()
	sLog := slog.New(slog.NewTextHandler(f, nil))
	s, err := server.NewServer(Hostname, Port, dbPath, sLog)
	if err != nil {
		panic(err)
	}
	if err := s.ListenAndServe(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
