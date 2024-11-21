package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikarwacki/gator/internal/database"
)

func users(st *state, c command) error {
	if len(c.Args) != 0 {
		return errors.New("Command doesn't take any arguments")
	}

	users, err := st.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == st.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}

	return nil
}

func resetUsers(st *state, c command) error {
	if len(c.Args) != 0 {
		return errors.New("Command doesn't take any arguments")
	}
	err := st.db.Reset(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Reset of users table was successful")

	return nil
}

func registerUser(st *state, c command) error {
	if len(c.Args) != 1 {
		return errors.New("Provide the name of the registered user")
	}

	name := c.Args[0]
	dt := time.Now().UTC()
	id := uuid.New()
	queryParams := database.CreateUserParams{ID: id, CreatedAt: dt, UpdatedAt: dt, Name: name}
	user, err := st.db.CreateUser(context.Background(), queryParams)
	if err != nil {
		return err
	}
	err = st.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully")
	printUser(user)

	return nil
}

func handlerLogin(st *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("Provide username for login command")
	}
	username := cmd.Args[0]
	user, err := st.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	err = st.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User has been successfully set for username: %s\n", username)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
