package main

import (
	"context"
	"fmt"

	"github.com/mikarwacki/gator/internal/database"
)

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {

	return func(s *state, c command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Current user is not logged in %v", err)
		}

		err = handler(s, c, user)
		if err != nil {
			return err
		}
		return nil
	}
}
