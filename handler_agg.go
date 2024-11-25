package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mikarwacki/gator/internal/database"
	"github.com/mikarwacki/gator/internal/rss"
)

const url = "https://www.wagslane.dev/index.xml"

func browse(s *state, c command) error {
	limit := 2
	if len(c.Args) == 1 {
		temp, err := strconv.Atoi(c.Args[0])
		if err != nil {
			return fmt.Errorf("Command requires one argument that is the number %v", err)
		}
		limit = temp
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("User is not registered %v", err)
	}

	params := database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Posts fetching failed %v", err)
	}

	for _, post := range posts {
		printPost(post)
	}
	return nil
}

func agg(st *state, c command) error {
	if len(c.Args) != 1 {
		return errors.New("This command requieres one argument (duration between fetches)")
	}
	timeBetweenReqs := c.Args[0]
	duration, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("Cound't parse duration %v", err)
	}
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		err = scrapeFeeds(st)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Failed getting next feed from db\n %v", err)
	}

	markedFeed, err := s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("Failed marking feed as fetched\n %v", err)
	}

	rssFeed, err := rss.FetchFeed(context.Background(), markedFeed.Url)
	if err != nil {
		return fmt.Errorf("Failed fetching the feed\n %v", err)
	}

	for _, rssItem := range rssFeed.Channel.Item {
		id := uuid.New()
		dt := time.Now().UTC()
		date, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", rssItem.PubDate)
		if err != nil {
			return fmt.Errorf("Coudn't parse Publish date of a post %v", err)
		}
		fmt.Println("____DEEBUUUUUG_____")
		fmt.Println(rssItem.Description)
		params := database.CreatePostParams{
			ID:          id,
			CreatedAt:   dt,
			UpdatedAt:   dt,
			Url:         rssItem.Link,
			FeedID:      markedFeed.ID,
			Title:       rssItem.Title,
			Description: {rssItem.Description,
			PostedAt:    date,
		}

		_, err = s.db.CreatePost(context.Background(), params)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			} else {
				log.Printf("Couldn't create post: %v", err)
				continue
			}
		}
	}
	return nil
}

func printPost(post database.GetPostsForUserRow) {
	fmt.Printf("%s from %s\n", post.PostedAt.Time.Format("Mon Jan 2"), post.FeedName)
	fmt.Printf("--- %s ---\n", post.Title)
	fmt.Printf("    %v\n", post.Description)
	fmt.Printf("Link: %s\n", post.Url)
	fmt.Println("=====================================")
}
