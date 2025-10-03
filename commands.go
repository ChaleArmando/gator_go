package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ChaleArmando/gator_go/internal/config"
	"github.com/ChaleArmando/gator_go/internal/database"
	"github.com/google/uuid"
)

type state struct {
	conf      *config.Config
	dbQueries *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login expect a single argument: username")
	}

	i, err := s.dbQueries.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	err = s.conf.SetUser(i.Name)
	if err != nil {
		return fmt.Errorf("user set failed: %w", err)
	}
	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("register expect a single argument: name")
	}

	dbArgs := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	i, err := s.dbQueries.CreateUser(context.Background(), dbArgs)
	if err != nil {
		return fmt.Errorf("create user failed: %w", err)
	}
	err = s.conf.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("user set failed: %w", err)
	}
	fmt.Println("User was created")
	printUser(i)
	return nil
}

func handlerReset(s *state, _ command) error {
	err := s.dbQueries.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("reset users failed: %w", err)
	}
	fmt.Println("Reset users was successful")
	return nil
}

func handlerUsers(s *state, _ command) error {
	users, err := s.dbQueries.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}
	for _, user := range users {
		status := ""
		if s.conf.CurrentUserName == user.Name {
			status = " (current)"
		}
		fmt.Println(user.Name + status)
	}
	return nil
}

func printUser(i database.User) {
	fmt.Printf("ID: %v\n", i.ID)
	fmt.Printf("Name: %v\n", i.Name)
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.cmds[cmd.name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
