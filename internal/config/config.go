package config

import (
	"io"
	"log/slog"

	"github.com/Blustak/go-transActor/internal/database"
)

type Config struct {
	Log    *slog.Logger
	DB     *database.Queries
	Reader io.Reader
	Writer io.Writer
}
