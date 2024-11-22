package main

import (
	"errors"
	"fmt"

	"github.com/mikarwacki/gator/internal/config"
	"github.com/mikarwacki/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.registeredCommands[cmd.Name]; ok {
		err := f(s, cmd)
		if err != nil {
			fmt.Println("Error running the command")
			return err
		}
		return nil
	}
	return errors.New("Command is not available")
}
