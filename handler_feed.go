package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shavits/boot-gator/internal/database"
)



func handlerAddFeed(s *state, cmd command, user database.User) error{
	if len(cmd.args) < 2{
		return fmt.Errorf("args empty, valid name and url expected")
	}
	name := cmd.args[0]
	url := cmd.args[1]

	params := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil{
		return fmt.Errorf("failed to create feed - %s", err)
	}
	printFeed(feed)


	
	followParams := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), followParams)
	if err != nil{
		return fmt.Errorf("failed to create feed_follow - %v", err)
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error{
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil{
		return fmt.Errorf("failed to get feeds - %s", err)
	}



	for _, feed := range(feeds){
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil{
		return fmt.Errorf("error getting user - %s", s.config.CurrentUserName)
		}
		fmt.Printf("name - %s\n", feed.Name)
		fmt.Printf("url - %s\n", feed.Name)
		fmt.Printf("user - %s\n", user.Name)
	}
	return nil
}


func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:      %v\n", feed.ID)
	fmt.Printf(" * Name:    %v\n", feed.Name)
	fmt.Printf(" * URL:    %v\n", feed.Url)
	fmt.Printf(" * USER_ID:    %v\n", feed.UserID)
}