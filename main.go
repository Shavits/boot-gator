package main

import (
	"fmt"
	"os"

	"github.com/shavits/boot-gator/internal/config"
)

type state struct{
	config *config.Config
}

func main() {
	curConfig, err := config.Read()
	if err != nil {
		fmt.Print(err)
	}
	curState := state{
		config: &curConfig,
	}
	cmds := Commands{handlers: make(map[string]func(*state, Command) error)}
	cmds.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2{
		fmt.Println("No arguments provided")
		os.Exit(1)
	}
	cmd := Command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&curState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
}

func handlerLogin(s *state, cmd Command) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("args empty, valid username expected")
	}
	user := cmd.args[0]
	err := s.config.SetUser(user)
	if err != nil{
		return fmt.Errorf("unable to set user - %s", err)
	}
	
	fmt.Printf("User set as - %s\n", user)
	return nil
}


