package utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/config"
	"log"
	"time"
)

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

func ParseTime(feed config.RSSItem) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05.000000", feed.PubDate)
	if err != nil {
		log.Printf("HandleAggregation error: unable to parse time %v\n", err)
		return time.Time{}, err
	}
	return parsedTime, nil
}
