package middleware

import (
	"context"
	"gator/internal/commands"
	"gator/internal/database"
	"gator/internal/state"
)

type HandlerFunc func(s *state.State, cmd commands.Command, user database.User) error

func LoggedIn(handler HandlerFunc) func(s *state.State, cmd commands.Command) error {
	return func(s *state.State, cmd commands.Command) error {
		ctx := context.Background()
		user, err := s.Db.GetUser(ctx, s.Config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
