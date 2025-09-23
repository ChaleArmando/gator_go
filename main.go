package main

import (
	"fmt"
	"os"

	"github.com/ChaleArmando/gator_go/internal/config"
)

func main() {
	conf := config.Read()

	newState := state{
		conf: &conf,
	}

	gatorCommands := commands{
		cmds: map[string]func(*state, command) error{},
	}

	gatorCommands.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("arguments not given")
		os.Exit(1)
	}
	cmdName := args[1]
	cmdArgs := args[2:]

	cmdInstance := command{name: cmdName, args: cmdArgs}
	err := gatorCommands.run(&newState, cmdInstance)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(config.Read())
}
