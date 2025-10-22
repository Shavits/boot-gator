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
	fmt.Printf("Create user - %s\n", userName)
	return nil

}