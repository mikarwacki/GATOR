package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
		params := database.CreatePostParams{
			ID:          id,
			CreatedAt:   dt,
			UpdatedAt:   dt,
			Url:         markedFeed.Url,
			FeedID:      markedFeed.ID,
			Title:       rssItem.Title,
			Description: rssItem.Description,
			PostedAt:    date,
		}

		_, err = s.db.CreatePost(context.Background(), params)
		if err != nil && err.(*pq.Error).Code != "23505" { //&& err = UrlFetchError
			fmt.Println(params.FeedID)
			fmt.Println(params.Url)
			return fmt.Errorf("Creating post error %v", err)
		}
	}
	return nil
}

func printRssItem(item rss.RSSItem) {
	fmt.Printf("Title: %v, Link: %v, Desc: %v, PubDate: %v\n", item.Title, item.Link, item.Description, item.PubDate)
}

func printPost(post database.Post) {
	fmt.Printf("Title: %v\n", post.Title)
	fmt.Printf("Url: %v\n", post.Url)
	fmt.Printf("Publication Date: %v\n", post.PostedAt)
	fmt.Printf("Description: %v\n", post.Description)
}
