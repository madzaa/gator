package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/commands"
	"gator/internal/database"
	"gator/internal/state"
	"time"

	"github.com/google/uuid"
)

func HandleLogin(s *state.State, cmd commands.Command) error {
	args := cmd.Arguments
	if len(args) == 0 || len(args) > 1 {
		return errors.New("invalid username")
	}
	username := args[0]

	exists, err := checkUserExists(s, username)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user does not exist")
	}

	err = s.Config.SetUser(username)
	if err != nil {
		return fmt.Errorf("unable to set user %v: %v ", username, err)
	}
	fmt.Printf("user has been set to %s", username)
	return nil
}

func HandleRegistration(s *state.State, cmd commands.Command) error {
	args := cmd.Arguments
	if len(args) == 0 || len(args) > 1 {
		return errors.New("invalid username")
	}
	username := args[0]

	exists, err := checkUserExists(s, username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	_, err = s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("unable to create user %v: %v", username, err)
	}

	fmt.Printf("user has been created %s", username)

	err = s.Config.SetUser(username)

	if err != nil {
		return fmt.Errorf("unable to set user %v: %v", username, err)
	}
	return nil
}

func HandleDeletion(s *state.State, cmd commands.Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to delete users: %v", err)
	}
	return nil
}

func HandleListUsers(s *state.State, cmd commands.Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to list users: %v", err)
	}
	for _, user := range users {
		if user == s.Config.CurrentUserName {
			fmt.Printf("* %v (current)", user)
		} else {
			fmt.Printf("* %v", user)
		}

	}
	return nil
}

func HandleAggregation(*state.State, commands.Command) error {
	feed, err := commands.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func HandleAddFeed(s *state.State, cmd commands.Command, user database.User) error {
	ctx := context.Background()
	args := cmd.Arguments
	feed, err := s.Db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      args[0],
		Url:       args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	_, err = s.Db.CreateFeedFollows(ctx, database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return err
	}

	fetchFeed, err := commands.FetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}
	fmt.Println(fetchFeed)

	return nil

}

func HandleListFeeds(s *state.State, cmd commands.Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("%s\n%s\n", feed.Name, feed.Name_2.String)
	}

	return nil
}

func HandleFeedFollow(s *state.State, cmd commands.Command, user database.User) error {
	ctx := context.Background()
	args := cmd.Arguments
	feed, err := s.Db.GetFeed(ctx, args[0])
	if err != nil {
		return err
	}
	_, err = s.Db.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println(feed.Name, user.Name)
	return nil
}

func HandleFeedFollowing(s *state.State, cmd commands.Command, user database.User) error {
	ctx := context.Background()
	feeds, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Println(feed.Name_2.String)
	}
	return nil
}

func checkUserExists(s *state.State, username string) (bool, error) {
	_, err := s.Db.GetUser(context.Background(), username)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, fmt.Errorf("query error: %w", err)
}
