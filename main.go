package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mikarwacki/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	s := state{cfg: &cfg}
	commands := commands{registeredCommands: make(map[string]func(*state, command) error)}
	commands.registeredCommands["login"] = handlerLogin
	ar := os.Args
	if len(ar) < 2 {
		log.Fatal("Not enough arguments provided")
	}

	command := command{Name: ar[1], Args: ar[2:]}
	fmt.Println(command.Name)
	fmt.Println(command.Args)
	err = commands.run(&s, command)
	if err != nil {
		log.Fatal(err)
	}
}
