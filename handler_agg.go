package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/shavits/boot-gator/internal/database"
)


func handlerAgg(s *state, cmd command) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("args empty, valid duration expected")
	}
	time_between_reqs , err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to parse duration - %s", cmd.args[0])
	}

	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}


func scrapeFeeds(s *state){
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil{
		fmt.Printf("error getting next feed - %v\n", err)
	}

	err = s.db.MarkFeedFetched(context.Background(),nextFeed.ID)
	if err != nil{
		fmt.Printf("error marking next feed as fetched - %v\n", err)
	}

	feedRss, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil{
		fmt.Printf("error fetching feed - %v\n", err)
	}

	//fmt.Printf("\n\n#####items for %s:\n", feedRss.Channel.Title)
	for _, item := range(feedRss.Channel.Item){
		//fmt.Println(item.Title)
		publishedAt, err := parsePubDate(item.PubDate)
        if err != nil {
            fmt.Printf("failed to parse time - %s: %v\n", item.PubDate, err)
            continue
        }
		params := database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
			Url: item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
			FeedID: nextFeed.ID,
		}

        _, err = s.db.CreatePost(context.Background(), params)
        if err != nil {
            // Ignore duplicate-url unique constraint (Postgres error code 23505).
            var pgErr *pgconn.PgError
            if errors.As(err, &pgErr) && pgErr.Code == "23505" {
                continue
            }
            var pqErr *pq.Error
            if errors.As(err, &pqErr) && string(pqErr.Code) == "23505" {
                continue
            }

            fmt.Printf("unable to create post - %v\n", err)
        }
	}
}


func parsePubDate(s string) (time.Time, error) {
    s = strings.TrimSpace(s)
    layouts := []string{
        time.RFC1123Z,                // "Mon, 02 Jan 2006 15:04:05 -0700"
        time.RFC1123,                 // "Mon, 02 Jan 2006 15:04:05 MST"
        time.RFC822Z,                 // "02 Jan 06 15:04 -0700"
        time.RFC822,                  // "02 Jan 06 15:04 MST"
        time.RFC3339,                 // "2006-01-02T15:04:05Z07:00"
        time.RFC3339Nano,
        "Mon, 02 Jan 2006 15:04:05 -0700",
        "02 Jan 2006 15:04:05 -0700",
        "Mon, 02 Jan 2006 15:04:05 MST",
        "02 Jan 2006 15:04:05 MST",
        time.ANSIC,
        time.UnixDate,
        time.RubyDate,
    }
    for _, l := range layouts {
        if t, err := time.Parse(l, s); err == nil {
            return t, nil
        }
    }
    return time.Time{}, fmt.Errorf("unrecognized time format: %q", s)
}
