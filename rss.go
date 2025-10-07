package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
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
	rss, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(rss)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("login expect arguments: feed name and feed url")
	}

	user, err := s.dbQueries.GetUser(context.Background(), s.conf.CurrentUserName)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	dbArgs := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}
	feed, err := s.dbQueries.CreateFeed(context.Background(), dbArgs)
	if err != nil {
		return fmt.Errorf("create user failed: %w", err)
	}

	printFeed(feed)
	return nil
}

func printFeed(i database.Feed) {
	fmt.Printf("ID: %v\n", i.ID)
	fmt.Printf("Name: %v\n", i.Name)
	fmt.Printf("URL: %v\n", i.Url)
	fmt.Printf("User ID: %v\n", i.UserID)
}
