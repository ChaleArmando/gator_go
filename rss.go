package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ChaleArmando/gator_go/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var rss RSSFeed
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &rss, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &rss, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &rss, err
	}

	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return &rss, err
	}

	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for i, val := range rss.Channel.Item {
		rss.Channel.Item[i].Title = html.UnescapeString(val.Title)
		rss.Channel.Item[i].Description = html.UnescapeString(val.Description)
	}

	return &rss, nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("login expect argument: time between requests")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't use argument for time requests: %w", err)
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.dbQueries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feed: %w", err)
	}

	rss, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	dbArgs := database.MarkFeedFetchedParams{
		ID:        feed.ID,
		UpdatedAt: time.Now(),
	}
	err = s.dbQueries.MarkFeedFetched(context.Background(), dbArgs)
	if err != nil {
		return fmt.Errorf("failed to mark fetched feed: %w", err)
	}

	for _, rssItem := range rss.Channel.Item {

		dbPostArgs := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       rssItem.Title,
			Url:         rssItem.Link,
			Description: database.StringNull(rssItem.Description),
			FeedID:      feed.ID,
			PublishedAt: database.TimeNull(rssItem.PubDate),
		}
		_, err = s.dbQueries.CreatePost(context.Background(), dbPostArgs)
		if err != nil {
			if !strings.Contains(err.Error(), "unique constraint \"posts_url_key\"") {
				log.Println(err)
			}
		}

	}
	return nil
}
