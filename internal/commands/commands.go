package commands

import (
	"context"
	"fmt"
	"gator/internal/config"
	"log"
)

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	handlers map[string]func(context.Context, *config.State, Command) error
}

func New() *Commands {
	return &Commands{handlers: make(map[string]func(context.Context, *config.State, Command) error)}
}

func (c *Commands) Run(ctx context.Context, s *config.State, cmd Command) error {
	if c.handlers[cmd.Name] == nil {
		err := fmt.Errorf("command %v does not exist", cmd.Name)
		log.Printf("Commands.Run error: %v\n", err)
		return err
	}
	err := c.handlers[cmd.Name](ctx, s, cmd)
	if err != nil {
		log.Printf("Commands.Run error: unable to Run Command %v: %v\n", cmd.Name, err)
		return fmt.Errorf("unable to Run Command: %v\n", err)
	}
	return err
}

func (c *Commands) Register(name string, f func(context.Context, *config.State, Command) error) {
	c.handlers[name] = f
}
