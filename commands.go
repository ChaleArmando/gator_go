package main

import (
	"errors"
	"fmt"

	"github.com/ChaleArmando/gator_go/internal/config"
)

type state struct {
	conf *config.Config
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
	err := s.conf.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
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
