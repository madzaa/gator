package main

import (
	"database/sql"
	"gator/internal/commands"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/handlers"
	"gator/internal/state"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("unable to read config: %v", err)
	}

	cfgState := &state.State{Config: &cfg}
	dbUrl := cfgState.Config.DbUrl
	db, err := sql.Open("postgres", dbUrl)
	cfgState.Db = database.New(db)

	cmds := commands.New()
	handlers.RegisterCommands(cmds)

	args := os.Args
	cmd := commands.Command{
		Name:      args[1],
		Arguments: args[2:],
	}
	err = cmds.Run(cfgState, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
