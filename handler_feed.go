package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikarwacki/gator/internal/database"
)

func addfeed(st *state, c command) error {
	if len(c.Args) != 2 {
		return errors.New("Command requires two arguments")
	}

	name := c.Args[0]
	url := c.Args[1]
	dt := time.Now().UTC()
	id := uuid.New()
	currentUser, err := st.db.GetUser(context.Background(), st.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feedParams := database.CreateFeedParams{ID: id, CreatedAt: dt, UpdatedAt: dt, Name: name, Url: url, UserID: currentUser.ID}
	feed, err := st.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}
	fmt.Println("Feed created successfully")
	printFeed(feed)

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:      %v\n", feed.ID)
	fmt.Printf(" * Name:    %v\n", feed.Name)
	fmt.Printf(" * Url:     %v\n", feed.Url)
	fmt.Printf(" * UserId:  %v\n", feed.UserID)
}
