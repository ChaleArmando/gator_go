package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ChaleArmando/gator_go/internal/config"
	"github.com/ChaleArmando/gator_go/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	conf      *config.Config
	dbQueries *database.Queries
}

func main() {
	conf := config.Read()

	db, err := sql.Open("postgres", conf.DbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)

	newState := state{
		conf:      &conf,
		dbQueries: dbQueries,
	}

	gatorCommands := commands{
		cmds: map[string]func(*state, command) error{},
	}

	gatorCommands.register("login", handlerLogin)
	gatorCommands.register("register", handlerRegister)
	gatorCommands.register("reset", handlerReset)
	gatorCommands.register("users", handlerUsers)
	gatorCommands.register("agg", handlerAgg)
	gatorCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	gatorCommands.register("feeds", handlerFeeds)
	gatorCommands.register("follow", middlewareLoggedIn(handlerFollow))
	gatorCommands.register("following", handlerFollowing)
	gatorCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	args := os.Args
	if len(args) < 2 {
		fmt.Println("arguments not given")
		os.Exit(1)
	}
	cmdName := args[1]
	cmdArgs := args[2:]

	cmdInstance := command{name: cmdName, args: cmdArgs}
	err = gatorCommands.run(&newState, cmdInstance)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//fmt.Println(config.Read())
}
