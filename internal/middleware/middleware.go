package middleware

import (
	"context"
	"gator/internal/commands"
	"gator/internal/config"
	"gator/internal/database"
	"log"
)

type HandlerFunc func(ctx context.Context, s *config.State, cmd commands.Command, user database.User) error

func LoggedIn(handler HandlerFunc) func(ctx context.Context, s *config.State, cmd commands.Command) error {
	return func(ctx context.Context, s *config.State, cmd commands.Command) error {
		user, err := s.Db.GetUser(ctx, s.Config.CurrentUserName)
		if err != nil {
			log.Printf("LoggedIn error: failed to get user %s: %v\n", s.Config.CurrentUserName, err)
			return err
		}
		return handler(ctx, s, cmd, user)
	}
}
