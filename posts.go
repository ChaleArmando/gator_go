package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ChaleArmando/gator_go/internal/database"
)

func handlerBrowse(s *state, cmd command) error {
	limit := "2"
	if len(cmd.args) == 1 {
		limit = cmd.args[0]
	} else if len(cmd.args) > 1 {
		fmt.Print(cmd.args)
		return errors.New("too many arguments given, browse expect argument: limit")
	}
	limitNum, err := strconv.Atoi(limit)
	if err != nil {
		return fmt.Errorf("argument given of incorrect type, not a number: %w", err)
	}
	posts, err := s.dbQueries.GetPostsForUser(context.Background(), int32(limitNum))
	if err != nil {
		return fmt.Errorf("posts not found: %w", err)
	}

	for _, post := range posts {
		feed, err := s.dbQueries.GetFeedByID(context.Background(), post.FeedID)
		if err != nil {
			return fmt.Errorf("feed not found: %w", err)
		}

		printPost(post, feed)
		fmt.Println("----------------------------------------------")
	}
	return nil
}

func printPost(post database.Post, feed database.Feed) {
	fmt.Printf("Title: %v\n", post.Title)
	fmt.Printf("Description: %v\n", post.Description)
	fmt.Printf("URL: %v\n", post.Url)
	fmt.Printf("Published At: %v\n", post.PublishedAt)
	fmt.Printf("Feed: %v\n", feed.Name)
}
