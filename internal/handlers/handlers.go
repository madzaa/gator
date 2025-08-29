package handlers

import (
	"context"
	"errors"
	"fmt"
	"gator/internal/commands"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/service"
	"gator/internal/utils"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func HandleLogin(ctx context.Context, s *config.State, c commands.Command) error {
	args := c.Arguments
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

func HandleRegistration(ctx context.Context, s *config.State, c commands.Command) error {
	args := c.Arguments
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

func HandleDeletion(ctx context.Context, s *config.State, c commands.Command) error {
	err := s.Db.DeleteUsers(ctx)
	if err != nil {
		log.Printf("HandleDeletion error: unable to delete users: %v\n", err)
		return fmt.Errorf("unable to delete users: %v\n", err)
	}
	return nil
}

func HandleListUsers(ctx context.Context, s *config.State, c commands.Command) error {
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

func HandleAggregation(ctx context.Context, s *config.State, c commands.Command) error {
	args := c.Arguments
	duration := args[0]
	parsedDuration, err := time.ParseDuration(duration)
	if err != nil {
		log.Printf("HandleAggregation error: unable to parse duration %v\n", err)
		return err
	}

	ticker := time.NewTicker(parsedDuration)
	for ; ; <-ticker.C {
		feeds, err := service.ScrapeFeed(ctx, s)
		if err != nil {
			log.Printf("HandleAggregation error: unable to scrape feeds %v\n", err)
			return err
		}
		fetchedFeed, err := service.FetchFeed(ctx, feeds.Url)
		if err != nil {
			log.Printf("HandleAggregation error: unable to fetch feeds %v\n", err)
			return err
		}
		err = service.CreatePost(s, fetchedFeed, ctx, feeds)
		if err != nil {
			log.Printf("HandleAggregation error: unable to create post %v\n", err)
			return err
		}
	}
}

func HandleAddFeed(ctx context.Context, s *config.State, c commands.Command, user database.User) error {
	args := c.Arguments
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

	fetchFeed, err := service.FetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}
	fmt.Println(fetchFeed.Channel.Title)

	return nil

}

func HandleListFeeds(ctx context.Context, s *config.State, c commands.Command) error {
	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for i, feed := range feeds {
		if i == 0 || feed.Username.String != feeds[i-1].Username.String {
			fmt.Printf("* user: %s\n", feed.Username.String)
		}
		fmt.Printf("  * feed: %s\n", feed.Feedname)
	}
	return nil
}

func HandleFeedFollow(ctx context.Context, s *config.State, c commands.Command, user database.User) error {
	args := c.Arguments
	feed, err := s.Db.GetFeedByUrl(ctx, args[0])
	if err != nil {
		return err
	}

	exists := false
	feedFollows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	for _, feedFollow := range feedFollows {
		if feedFollow.Feedname.String == feed.Url {
			exists = true
		}
	}

	if !exists {
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
		fmt.Printf("Created feed %s for user %s", feed.Name, user.Name)
	} else {
		fmt.Printf("Feed follow exists\n")
	}

	return nil
}

func HandleFeedFollowing(ctx context.Context, s *config.State, c commands.Command, user database.User) error {
	feeds, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Println(feed.Feedname.String)
	}
	return nil
}

func HandleUnfollow(ctx context.Context, s *config.State, c commands.Command, user database.User) error {
	args := c.Arguments
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

func HandleBrowse(ctx context.Context, s *config.State, c commands.Command, user database.User) error {
	args := c.Arguments
	limit := 1
	var err error
	if len(args) > 0 {
		limit, err = strconv.Atoi(args[0])
	}

	if err != nil {
		return err
	}

	fmt.Println(user.ID, limit)
	posts, err := s.Db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	fmt.Println(posts)
	return nil
}
