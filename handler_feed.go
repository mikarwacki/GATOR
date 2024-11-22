package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikarwacki/gator/internal/database"
)

func feeds(st *state, c command) error {
	if len(c.Args) != 0 {
		return errors.New("This command doesn't accept arguments")
	}

	feeds, err := st.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := st.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Coudn't find user %v", feed.UserID)
		}
		printFeed(feed, user)
	}
	return nil
}

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

	feedFollow := database.CreateFeedFollowParams{CreatedAt: dt, UpdatedAt: dt, UserID: currentUser.ID, FeedID: feed.ID}
	_, err = st.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return err
	}

	fmt.Println("Feed created successfully and followed")
	printFeed(feed, currentUser)

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf(" * ID:        %v\n", feed.ID)
	fmt.Printf(" * CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf(" * UpdatedAt: %v\n", feed.UpdatedAt)
	fmt.Printf(" * Name:      %v\n", feed.Name)
	fmt.Printf(" * Url:       %v\n", feed.Url)
	fmt.Printf(" * User:      %v\n", user.Name)
}
