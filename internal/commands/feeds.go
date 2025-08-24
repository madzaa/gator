package commands

import (
	"context"
	"encoding/xml"
	"gator/internal/config"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedUrl string) (*config.RSSFeed, error) {
	client := http.Client{}
	request, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "gator")

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rss config.RSSFeed
	body, err := io.ReadAll(resp.Body)
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, err
	}

	sanitiseFields(&rss)
	return &rss, err
}

func sanitiseFields(rss *config.RSSFeed) {
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
}
