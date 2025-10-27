package main

import (
	"context"
	"database/sql"
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
	if err == sql.ErrNoRows{
		return fmt.Errorf("user %s does not exist", s.config.CurrentUserName)
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


func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:      %v\n", feed.ID)
	fmt.Printf(" * Name:    %v\n", feed.Name)
	fmt.Printf(" * URL:    %v\n", feed.Url)
	fmt.Printf(" * USER_ID:    %v\n", feed.UserID)
}