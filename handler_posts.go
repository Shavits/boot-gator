package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/shavits/boot-gator/internal/database"
)



func handlerBrowse(s *state, cmd command, user database.User) error{
	limit := 2
    if len(cmd.args) >= 1 {
        parsed, err := strconv.Atoi(cmd.args[0])
        if err != nil {
            return fmt.Errorf("failed to parse limit - %v - %v", cmd.args[0], err)
        }
        if parsed > 0 {
            limit = parsed
        }
	}

	posts, err := s.db.GetPostsForUser(context.Background(), user.ID)
	if err != nil{
		return fmt.Errorf("failed to get posts for user %v - %v", user.ID, err)
	}

	for i := range(limit){
		printPost(posts[i])
	}
	return nil


}

func printPost(p database.Post) {
    pub := p.PublishedAt
    pubStr := "unknown"
    if !pub.IsZero() {
        pubStr = pub.Format(time.RFC1123)
    }
    fmt.Printf("Title: %s\nURL:   %s\nPublished: %s\n\n%s\n\n",
        p.Title, p.Url, pubStr, p.Description)
}