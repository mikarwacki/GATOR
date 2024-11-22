package main

import (
	"context"
	"fmt"

	"github.com/mikarwacki/gator/internal/rrs"
)

const url = "https://www.wagslane.dev/index.xml"

func agg(st *state, c command) error {
	rssFeed, err := rrs.FetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Println(rssFeed)
	return nil
}
