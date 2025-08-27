package config

import (
	"gator/internal/database"
)

type State struct {
	Db     *database.Queries
	Config *Config
}
