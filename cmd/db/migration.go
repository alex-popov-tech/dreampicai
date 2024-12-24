package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"dreampicai/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: [up|status|reset|create]")
		os.Exit(1)
	}

	command := os.Args[1]

	// Parse remaining flags
	flag.Parse()

	env, err := utils.ValidateEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("pgx", env.DatabaseDirectURL)
	if err != nil {
		log.Fatal(err)
	}

	switch command {
	case "up":
		err = goose.Up(db, env.DatabaseMigrations)
		if err != nil {
			log.Fatal(err)
		}
	case "down":
		err = goose.Up(db, env.DatabaseMigrations)
		if err != nil {
			log.Fatal(err)
		}
	case "status":
		err = goose.Status(db, env.DatabaseMigrations)
		if err != nil {
			log.Fatal(err)
		}
	case "reset":
		err = goose.Reset(db, env.DatabaseMigrations)
		if err != nil {
			log.Fatal(err)
		}
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Usage: create [name]")
			os.Exit(1)
		}
		name := os.Args[2]
		err = goose.Create(db, env.DatabaseMigrations, name, "sql")
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: [up|status|reset]")
		os.Exit(1)
	}
}
