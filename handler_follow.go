package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mikarwacki/gator/internal/database"
)

func unfollow(st *state, c command, user database.User) error {
	if len(c.Args) != 1 {
		return errors.New("This command requires exactly one argument")
	}
	url := c.Args[0]
	params := database.DeleteFeedByUserIdAndUrlParams{UserID: user.ID, Url: url}
	err := st.db.DeleteFeedByUserIdAndUrl(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("Feed %v has been deleted\n", url)
	return nil
}

func following(st *state, c command, user database.User) error {
	if len(c.Args) != 0 {
		return errors.New("This command doesn't accept arguments")
	}

	feedJoin, err := st.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Printing all followed feeds for user %v", user.Name)
	for _, row := range feedJoin {
		fmt.Printf("Feed name: %v\n", row.FeedName)
	}

	return nil
}

func follow(st *state, c command, user database.User) error {
	if len(c.Args) != 1 {
		return errors.New("This command requires exactly one argument")
	}

	url := c.Args[0]
	feed, err := st.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	dt := time.Now().UTC()
	followParams := database.CreateFeedFollowParams{CreatedAt: dt, UpdatedAt: dt, FeedID: feed.ID, UserID: user.ID}
	_, err = st.db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return err
	}

	printFollow(feed, user)
	return nil
}

func printFollow(feed database.Feed, user database.User) {
	fmt.Printf("feed name: %v", feed.Name)
	fmt.Printf("user name: %v", user.Name)
}
