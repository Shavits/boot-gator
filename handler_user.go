package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shavits/boot-gator/internal/database"
)	


func handlerLogin(s *state, cmd command) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("args empty, valid username expected")
	}
	userName := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), userName)
	if err == sql.ErrNoRows{
		return fmt.Errorf("user %s does not exist", userName)
	}
	err = s.config.SetUser(userName)
	if err != nil{
		return fmt.Errorf("unable to set user - %s", err)
	}
	
	fmt.Printf("User set as - %s\n", userName)
	return nil
}

func handlerRegister(s *state, cmd command) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("args empty, valid username expected")
	}
	userName := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), userName)
	if err == nil && user.Name !=  ""{
		return fmt.Errorf("user already exists")

	}

	params := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: userName,


	}
	user, err = s.db.CreateUser(context.Background(), params)
	if err != nil{
		return fmt.Errorf("failed to create user - %s", err)
	} 

	err = s.config.SetUser(userName)
	if err != nil{
		return err
	}
	fmt.Printf("Created user - %s\n", userName)
	printUser(user)
	return nil

}


func handlerReset(s *state, cmd command) error{
	err := s.db.ResetUsers(context.Background())
	if err != nil{
		return fmt.Errorf("unable to reset users - %s", err)
	}

	return nil
}


func handlerUsers(s *state, cmd command) error{
	users, err := s.db.GetUsers(context.Background())
	if err != nil{
		return fmt.Errorf("failed to get users - %s", err)
	}

	for _, user := range(users){
		is_current := ""
		if user.Name == s.config.CurrentUserName{
			is_current = " (current)"
		}
		fmt.Printf("%s%s\n", user.Name, is_current)
	}
	return nil
}


func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}