package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/shavits/boot-gator/internal/config"
	"github.com/shavits/boot-gator/internal/database"
)

type state struct{
	db *database.Queries
	config *config.Config
}

func main() {
	curConfig, err := config.Read()
	if err != nil {
		fmt.Print(err)
	}

	db, err := sql.Open("postgres", curConfig.DbURL)
	if err != nil {
		fmt.Printf("error opening db - %s", err)
	}
	dbQueries := database.New(db)

	curState := state{
		db: dbQueries,
		config: &curConfig,
	}

	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))
	args := os.Args
	if len(args) < 2{
		fmt.Println("No arguments provided")
		os.Exit(1)
	}
	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&curState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}







