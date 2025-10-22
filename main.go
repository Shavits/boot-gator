package main

import (
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




