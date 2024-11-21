package main

import (
	"errors"
	"fmt"
)

func handlerLogin(st *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("Provide username for login command")
	}
	username := cmd.Args[0]
	err := st.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User has been successfully set for username: %s\n", username)
	return nil
}
