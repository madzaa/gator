package handlers

import (
	"context"
	"errors"
	"fmt"
	"gator/internal/commands"
	"gator/internal/database"
	"gator/internal/state"
	"gator/internal/utils"
	"log"
	"time"

	"github.com/google/uuid"
)

func HandleLogin(ctx context.Context, s *state.State, cmd commands.Command) error {
	args := cmd.Arguments
	if len(args) == 0 || len(args) > 1 {
		err := errors.New("invalid username")
		log.Printf("HandleLogin error: %v\n", err)
		return err
	}
	username := args[0]

	exists, err := utils.CheckExists(ctx, func(ctx context.Context) error {
		_, err := s.Db.GetUser(ctx, username)
		return err
	})
	if err != nil {
		log.Printf("HandleLogin error: %v\n", err)
		return err
	}
	if !exists {
		err := errors.New("user does not exist")
		log.Printf("HandleLogin error: %v\n", err)
		return err
	}

	err = s.Config.SetUser(username)
	if err != nil {
		log.Printf("HandleLogin error: unable to set user %v: %v\n", username, err)
		return fmt.Errorf("unable to set user %v: %v\n", username, err)
	}
	fmt.Printf("user has been set to %s\n", username)
	return nil
}

func HandleRegistration(ctx context.Context, s *state.State, cmd commands.Command) error {
	args := cmd.Arguments
	if len(args) == 0 || len(args) > 1 {
		err := errors.New("invalid username")
		log.Printf("HandleRegistration error: %v\n", err)
		return err
	}
	username := args[0]

	exists, err := utils.CheckExists(ctx, func(ctx context.Context) error {
		_, err := s.Db.GetUser(ctx, username)
		return err
	})
	if err != nil {
		log.Printf("HandleRegistration error: %v\n", err)
		return err
	}
	if exists {
		err := errors.New("user already exists")
		log.Printf("HandleRegistration error: %v\n", err)
		return err
	}

	_, err = s.Db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		log.Printf("HandleRegistration error: unable to create user %v: %v\n", username, err)
		return fmt.Errorf("unable to create user %v: %v\n", username, err)
	}

	fmt.Printf("user has been created %s\n", username)

	err = s.Config.SetUser(username)

	if err != nil {
		log.Printf("HandleRegistration error: unable to set user %v: %v\n", username, err)
		return fmt.Errorf("unable to set user %v: %v\n", username, err)
	}
	return nil
}

func HandleDeletion(ctx context.Context, s *state.State, cmd commands.Command) error {
	err := s.Db.DeleteUsers(ctx)
	if err != nil {
		log.Printf("HandleDeletion error: unable to delete users: %v\n", err)
		return fmt.Errorf("unable to delete users: %v\n", err)
	}
	return nil
}

func HandleListUsers(ctx context.Context, s *state.State, cmd commands.Command) error {
	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		log.Printf("HandleListUsers error: unable to list users: %v\n", err)
		return fmt.Errorf("unable to list users: %v\n", err)
	}
	for _, user := range users {
		if user == s.Config.CurrentUserName {
			fmt.Printf("* %v (current)\n", user)
		} else {
			fmt.Printf("* %v\n", user)
		}

	}
	return nil
}

func HandleAggregation(ctx context.Context, s *state.State, cmd commands.Command) error {
	args := cmd.Arguments
	duration := args[0]
	parsedDuration, err := time.ParseDuration(duration)
	if err != nil {
		log.Printf("HandleAggregation error: unable to parse duration %v\n", err)
		return err
	}

	ticker := time.NewTicker(parsedDuration)
	for ; ; <-ticker.C {
		feeds, err := utils.ScrapeFeed(ctx, s)
		if err != nil {
			log.Printf("HandleAggregation error: unable to scrape feeds %v\n", err)
			return err
		}
		fetchedFeed, err := commands.FetchFeed(ctx, feeds.Url)
		if err != nil {
			log.Printf("HandleAggregation error: unable to fetch feeds %v\n", err)
			return err
		}
		err = utils.CreatePost(s, fetchedFeed, ctx, feeds)
		if err != nil {
			log.Printf("HandleAggregation error: unable to create post %v\n", err)
			return err
		}
	}
}

func HandleAddFeed(ctx context.Context, s *state.State, cmd commands.Command, user database.User) error {
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

func HandleListFeeds(ctx context.Context, s *state.State, cmd commands.Command) error {
	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("%s\n%s\n", feed.Name, feed.Name_2.String)
	}

	return nil
}

func HandleFeedFollow(ctx context.Context, s *state.State, cmd commands.Command, user database.User) error {
	args := cmd.Arguments
	feed, err := s.Db.GetFeedByUrl(ctx, args[0])
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
	fmt.Println(feed.Name, user.Name)
	return nil
}

func HandleFeedFollowing(ctx context.Context, s *state.State, cmd commands.Command, user database.User) error {
	feeds, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Println(feed.Name_2.String)
	}
	return nil
}

func HandleUnfollow(ctx context.Context, s *state.State, cmd commands.Command, user database.User) error {
	args := cmd.Arguments
	feed, err := s.Db.GetFeedByUrl(ctx, args[0])
	if err != nil {
		return err
	}
	err = s.Db.DeleteFeedFollowsForUser(ctx, database.DeleteFeedFollowsForUserParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		return err
	}
	return nil
}
