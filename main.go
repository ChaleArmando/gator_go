package main

import (
	"fmt"

	"github.com/ChaleArmando/gator_go/internal/config"
)

func main() {
	conf := config.Read()
	err := conf.SetUser("Armando")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config.Read())
}
