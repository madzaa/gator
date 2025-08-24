package commands

import (
	"fmt"
	"gator/internal/state"
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
		return fmt.Errorf("command %v does not exist", cmd.Name)
	}
	err := c.handlers[cmd.Name](s, cmd)
	if err != nil {
		return fmt.Errorf("unable to Run Command: %v", err)
	}
	return err
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	c.handlers[name] = f
}
