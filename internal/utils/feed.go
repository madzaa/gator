package utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/state"
	"log"
	"time"

	"github.com/google/uuid"
)

func ScrapeFeed(ctx context.Context, s *state.State) (database.Feed, error) {
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

func CheckExists(ctx context.Context, lookUp func(ctx context.Context) error) (bool, error) {
	err := lookUp(ctx)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, fmt.Errorf("query error: %w", err)
}

func CreatePost(s *state.State, fetchedFeed *config.RSSFeed, ctx context.Context, feeds database.Feed) error {
	for _, feed := range fetchedFeed.Channel.Item {
		_, err := CheckExists(ctx, func(ctx context.Context) error {
			_, err := s.Db.GetPost(ctx, feed.Link)
			return err
		})
		if err != nil {
			return err
		}
		parsedTime, err := parseTime(feed)
		if err != nil {
			return err
		}

		err = s.Db.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     feed.Title,
			Url:       feed.Link,
			Description: sql.NullString{
				String: feed.Description,
				Valid:  true,
			},
			PublishedAt: sql.NullTime{
				Time:  parsedTime.Local(),
				Valid: true,
			},
			FeedID: feeds.ID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func parseTime(feed config.RSSItem) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05.000000", feed.PubDate)
	if err != nil {
		log.Printf("HandleAggregation error: unable to parse time %v\n", err)
		return time.Time{}, err
	}
	return parsedTime, nil
}
