package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/commands"
	"gator/internal/database"
	"gator/internal/state"
	"log"
	"time"

	"github.com/google/uuid"
)

func HandleLogin(s *state.State, cmd commands.Command) error {
	args := cmd.Arguments
	if len(args) == 0 || len(args) > 1 {
		err := errors.New("invalid username")
		log.Printf("HandleLogin error: %v\n", err)
		return err
	}
	username := args[0]

	exists, err := checkUserExists(s, username)
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

func HandleRegistration(s *state.State, cmd commands.Command) error {
	args := cmd.Arguments
	if len(args) == 0 || len(args) > 1 {
		err := errors.New("invalid username")
		log.Printf("HandleRegistration error: %v\n", err)
		return err
	}
	username := args[0]

	exists, err := checkUserExists(s, username)
	if err != nil {
		log.Printf("HandleRegistration error: %v\n", err)
		return err
	}
	if exists {
		err := errors.New("user already exists")
		log.Printf("HandleRegistration error: %v\n", err)
		return err
	}

	_, err = s.Db.CreateUser(context.Background(), database.CreateUserParams{
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

func HandleDeletion(s *state.State, cmd commands.Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		log.Printf("HandleDeletion error: unable to delete users: %v\n", err)
		return fmt.Errorf("unable to delete users: %v\n", err)
	}
	return nil
}

func HandleListUsers(s *state.State, cmd commands.Command) error {
	users, err := s.Db.GetUsers(context.Background())
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

func HandleAggregation(s *state.State, cmd commands.Command) error {
	ctx := context.Background()
	args := cmd.Arguments
	duration := args[0]
	parsedDuration, err := time.ParseDuration(duration)
	if err != nil {
		log.Printf("HandleAggregation error: unable to parse duration %v\n", err)
		return err
	}

	ticker := time.NewTicker(parsedDuration)
	for ; ; <-ticker.C {
		feeds, err := scrapeFeed(ctx, s)
		if err != nil {
			log.Printf("HandleAggregation error: unable to scrape feeds %v\n", err)
			return err
		}
		feed, err := commands.FetchFeed(ctx, feeds.Url)
		if err != nil {
			log.Printf("HandleAggregation error: unable to fetch feeds%v\n", err)
			return err
		}
		for _, item := range feed.Channel.Item {
			fmt.Println(item.Title)
		}
	}
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
	feed, err := s.Db.GetFeedByUrl(ctx, args[0])
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

func HandleUnfollow(s *state.State, cmd commands.Command, user database.User) error {
	ctx := context.Background()
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

func scrapeFeed(ctx context.Context, s *state.State) (database.Feed, error) {
	next, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return database.Feed{}, err
	}

	err = s.Db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		ID: next.ID,
	})
	if err != nil {
		return database.Feed{}, err
	}

	feed, err := s.Db.GetFeedByUrl(ctx, next.Url)
	if err != nil {
		return database.Feed{}, err
	}
	return feed, nil
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
