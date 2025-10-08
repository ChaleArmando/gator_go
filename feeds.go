package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ChaleArmando/gator_go/internal/database"
	"github.com/google/uuid"
)

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
		return fmt.Errorf("create feed failed: %w", err)
	}

	feedFollowArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.dbQueries.CreateFeedFollow(context.Background(), feedFollowArgs)
	if err != nil {
		return fmt.Errorf("create feed follow failed: %w", err)
	}

	fmt.Println("Feed created successfully")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("Feed followed successfully")
	printFeedFollow(feedFollow)
	return nil
}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.dbQueries.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %w", err)
	}

	for _, feed := range feeds {
		user, err := s.dbQueries.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		printFeed(feed, user)
		fmt.Println("----------------------------------------------")
	}
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("Created at: %v\n", feed.CreatedAt)
	fmt.Printf("Updated at: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("URL: %v\n", feed.Url)
	fmt.Printf("User Name: %v\n", user.Name)
}
