package utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/config"
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
	timeFormats := []string{
		time.RFC3339Nano,                  // "2006-01-02T15:04:05.999999999Z07:00"
		time.RFC3339,                      // "2006-01-02T15:04:05Z07:00"
		time.RFC1123Z,                     // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,                      // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC822Z,                      // "02 Jan 06 15:04 -0700"
		time.RFC822,                       // "02 Jan 06 15:04 MST"
		time.RFC850,                       // "Monday, 02-Jan-06 15:04:05 MST"
		time.ANSIC,                        // "Mon Jan _2 15:04:05 2006"
		time.UnixDate,                     // "Mon Jan _2 15:04:05 MST 2006"
		time.RubyDate,                     // "Mon Jan 02 15:04:05 -0700 2006"
		time.Kitchen,                      // "3:04PM"
		time.Stamp,                        // "Jan _2 15:04:05"
		time.StampMilli,                   // "Jan _2 15:04:05.000"
		time.StampMicro,                   // "Jan _2 15:04:05.000000"
		time.StampNano,                    // "Jan _2 15:04:05.000000000"
		"2006-01-02 15:04:05.000000",      // Your microsecond format
		"2006-01-02 15:04:05.000",         // Millisecond format
		"2006-01-02 15:04:05",             // Basic datetime
		"2006-01-02T15:04:05",             // ISO without timezone
		"2006-01-02",                      // Date only
		"01/02/2006 15:04:05",             // US format with time
		"01/02/2006",                      // US date format
		"02/01/2006 15:04:05",             // European format with time
		"02/01/2006",                      // European date format
		"2006/01/02 15:04:05",             // Alternative format
		"2006/01/02",                      // Alternative date format
		"Jan 2, 2006 3:04:05 PM",          // Readable format
		"Jan 2, 2006",                     // Readable date
		"January 2, 2006 3:04:05 PM",      // Full readable format
		"January 2, 2006",                 // Full readable date
		"Mon, 02 Jan 2006 15:04:05 -0700", // RFC1123 with numeric timezone
	}
	parsedTime := time.Time{}
	var err error
	for _, format := range timeFormats {
		parsedTime, err = time.Parse(format, feed.PubDate)
		if err == nil {
			break
		}
	}
	if err != nil {
		return time.Time{}, errors.New("no format found")

	}
	return parsedTime, nil
}
