package config

import (
	"log/slog"
	"sync"

	"github.com/Blustak/go-transActor/internal/database"
)

type State struct {
	sync.RWMutex
	*slog.Logger
	DB *database.Queries
}
