package main

import (
	"flag"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq" // Подключаем PostgreSQL драйвер
	"github.com/pressly/goose/v3"
)

const LocalFilePath = "config-yaml/local.yaml"
const migrationsDir = "./migrations"

func main() {
	dbDriver := "postgres"
	dbString := "postgres://postgres:root@localhost:5432/EventServiceDB?sslmode=disable"
	// open db
	db, err := goose.OpenDBWithDriver(dbDriver, dbString)
	if err != nil {
		log.Fatalf("failed to open database: %s: %v", dbString, err)
	}
	defer db.Close()

	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatalf("usage: migrator [command] [arguments]")
	}

	command := flag.Arg(0)

	switch command {
	case "up":
		err = goose.Up(db, migrationsDir)
	case "down":
		err = goose.Down(db, migrationsDir)
	default:
		log.Fatalf("unknow command %s", command)
	}
	if err != nil {
		log.Fatalf("failed to execute command: %s ,%v", command, err)
	}
	log.Printf("command %s execute successfully", command)
}
