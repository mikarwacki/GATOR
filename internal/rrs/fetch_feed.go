package rrs

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}

	return &rssFeed, nil
}
