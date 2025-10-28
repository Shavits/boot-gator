package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shavits/boot-gator/internal/database"
)



func handlerAddFeed(s *state, cmd command) error{
		if len(cmd.args) < 2{
		return fmt.Errorf("args empty, valid name and url expected")
	}
	name := cmd.args[0]
	url := cmd.args[1]
	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil{
		return fmt.Errorf("error getting user - %s", s.config.CurrentUserName)
	}

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