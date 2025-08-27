package service

import (
	"context"
	"database/sql"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/utils"
	"time"

	"github.com/google/uuid"
)

func CreatePost(s *config.State, fetchedFeed *config.RSSFeed, ctx context.Context, feeds database.Feed) error {
	for _, feed := range fetchedFeed.Channel.Item {
		_, err := utils.CheckExists(ctx, func(ctx context.Context) error {
			_, err := s.Db.GetPost(ctx, feed.Link)
			return err
		})
		if err != nil {
			return err
		}
		parsedTime, err := utils.ParseTime(feed)
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
