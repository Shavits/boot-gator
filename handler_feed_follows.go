package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shavits/boot-gator/internal/database"
)


func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0{
		return fmt.Errorf("args empty, valid url expected")
	}
	url := cmd.args[0]



	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil{
		return fmt.Errorf("error getting feed for url - %s", url)
	}

	params := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil{
		return fmt.Errorf("failed to create feed_follow - %v", err)
	}

	printFllow(follow)
	return nil


}


func handlerFollowing(s *state, cmd command, user database.User) error{
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil{
		return fmt.Errorf("error getting feeds - %v", err)
	}
	
	for _, follow := range(follows){
		fmt.Println(follow.FeedName)
	}
	return nil 
}



func printFllow(follow database.CreateFeedFollowRow) {
	fmt.Printf(" * User:      %v\n", follow.UserName)
	fmt.Printf(" * Feed:    %v\n", follow.FeedName)
}