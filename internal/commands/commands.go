package commands

import (
	"fmt"
	"gator/internal/state"
	"log"
)

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	handlers map[string]func(*state.State, Command) error
}

func New() *Commands {
	return &Commands{handlers: make(map[string]func(*state.State, Command) error)}
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	if c.handlers[cmd.Name] == nil {
		err := fmt.Errorf("command %v does not exist", cmd.Name)
		log.Printf("Commands.Run error: %v", err)
		return err
	}
	err := c.handlers[cmd.Name](s, cmd)
	if err != nil {
		log.Printf("Commands.Run error: unable to Run Command %v: %v", cmd.Name, err)
		return fmt.Errorf("unable to Run Command: %v", err)
	}
	return err
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	c.handlers[name] = f
}
