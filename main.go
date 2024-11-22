package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/mikarwacki/gator/internal/config"
	"github.com/mikarwacki/gator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal("Cound't connect to the db")
	}

	dbQueries := database.New(db)

	s := state{cfg: &cfg, db: dbQueries}
	commands := commands{registeredCommands: make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)
	commands.register("register", registerUser)
	commands.register("reset", resetUsers)
	commands.register("users", users)
	commands.register("agg", agg)
	commands.register("addfeed", addfeed)
	commands.register("feeds", feeds)
	ar := os.Args
	if len(ar) < 2 {
		log.Fatal("Not enough arguments provided")
	}

	command := command{Name: ar[1], Args: ar[2:]}
	err = commands.run(&s, command)
	if err != nil {
		log.Fatal(err)
	}
}
