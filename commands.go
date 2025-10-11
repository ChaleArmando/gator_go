package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/ChaleArmando/gator_go/internal/database"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.dbQueries.GetUser(context.Background(), s.conf.CurrentUserName)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		return handler(s, cmd, user)
	}
}
