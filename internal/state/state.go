package state

import (
	config "gator/internal/config"
	"gator/internal/database"
)

type State struct {
	Db     *database.Queries
	Config *config.Config
}
