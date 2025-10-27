package main

import (
	"context"
	"fmt"
)


func handlerAgg(s *state, cmd command) error{
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil{
		return err
	}
	fmt.Print(feed)
	return nil
}