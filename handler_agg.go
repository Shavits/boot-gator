package main

import (
	"context"
	"fmt"
	"time"
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

	fmt.Printf("\n\n#####items for %s:\n", feedRss.Channel.Title)
	for _, item := range(feedRss.Channel.Item){
		fmt.Println(item.Title)
	}
}