package service

import (
	"context"
	"database/sql"
	"encoding/xml"
	"gator/internal/config"
	"gator/internal/database"
	"html"
	"io"
	"log"
	"net/http"
	"time"
)

func FetchFeed(ctx context.Context, feedUrl string) (*config.RSSFeed, error) {
	client := http.Client{}
	request, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		log.Printf("FetchFeed error: failed to create request for %s: %v\n", feedUrl, err)
		return nil, err
	}
	request.Header.Set("User-Agent", "gator")

	resp, err := client.Do(request)
	if err != nil {
		log.Printf("FetchFeed error: failed to perform request for %s: %v\n", feedUrl, err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("FetchFeed error: failed to close response body for %s: %v\n", feedUrl, err)
			return
		}
	}(resp.Body)

	var rss config.RSSFeed
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("FetchFeed error: failed to read response body for %s: %v\n", feedUrl, err)
		return nil, err
	}
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		log.Printf("FetchFeed error: failed to unmarshal XML for %s: %v\n", feedUrl, err)
		return nil, err
	}

	sanitiseFields(&rss)
	return &rss, err
}

func ScrapeFeed(ctx context.Context, s *config.State) (database.Feed, error) {
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

func sanitiseFields(rss *config.RSSFeed) {
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
}
