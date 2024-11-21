package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikarwacki/gator/internal/database"
)

func registerUser(st *state, c command) error {
	if len(c.Args) != 1 {
		return errors.New("Provide the name of the registered user")
	}

	name := c.Args[0]
	dt := time.Now()
	id := uuid.NullUUID{UUID: uuid.New(), Valid: true}
	queryParams := database.CreateUserParams{ID: id, CreatedAt: dt, UpdatedAt: dt, Name: name}
	user, err := st.db.CreateUser(context.Background(), queryParams)
	if err != nil {
		return err
	}
	st.cfg.SetUser(user.Name)
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
