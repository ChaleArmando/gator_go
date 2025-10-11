package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ChaleArmando/gator_go/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("login expect arguments: feed url")
	}

	feed, err := s.dbQueries.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("feed not found: %w", err)
	}

	dbArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.dbQueries.CreateFeedFollow(context.Background(), dbArgs)
	if err != nil {
		return fmt.Errorf("create feed follow failed: %w", err)
	}

	fmt.Println("Feed followed successfully")
	printFeedFollow(feedFollow)
	return nil
}

func handlerFollowing(s *state, _ command) error {
	feedFollows, err := s.dbQueries.GetFeedFollowsForUser(context.Background(), s.conf.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get feed follows: %w", err)
	}

	fmt.Printf("Feeds followed by %v:\n", s.conf.CurrentUserName)
	for _, feedFollow := range feedFollows {
		fmt.Println(feedFollow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("login expect arguments: feed url")
	}

	feed, err := s.dbQueries.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("feed not found: %w", err)
	}

	dbArgs := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.dbQueries.DeleteFeedFollow(context.Background(), dbArgs)
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}

	fmt.Printf("%s unfollowed feed sucessfully\n", user.Name)
	fmt.Printf("Unfollowed Feed: %s\n", feed.Name)
	return nil
}

func printFeedFollow(feedFollow database.CreateFeedFollowRow) {
	fmt.Printf("User: %v\n", feedFollow.UserName)
	fmt.Printf("Feed: %v\n", feedFollow.FeedName)
}
