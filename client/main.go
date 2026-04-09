package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/Blustak/go-transActor/client/internal/client"
	"github.com/Blustak/go-transActor/internal/protocol"
	"github.com/joho/godotenv"
)

const Hostname = "localhost"
const Port = "8080"

func main() {
	godotenv.Load()
	protocol.Init()
	cli := client.Client{
		Hostname: Hostname,
		Port:     Port,
		In:       bufio.NewScanner(os.Stdin),
		Out:      bufio.NewWriter(os.Stdout),
	}
	f := os.Stdout
	var err error
	cLogPath, ok := os.LookupEnv("CLIENT_LOG")
	if ok {
		f, err = os.Open(cLogPath)
		if err != nil {
			panic(fmt.Sprintf("error opening log file: %v", err))
		}
	}
	defer f.Close()
	cli.Logger = slog.New(slog.NewTextHandler(f, nil))
	for {
		if err := cli.Run(); err != nil {
			cli.Error("error, quitting", slog.Any("error", err))
			break
		}
	}
}
