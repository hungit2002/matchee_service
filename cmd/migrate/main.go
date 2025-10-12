package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		direction string
		steps     int
	)
	flag.StringVar(&direction, "direction", "up", "up|down|goto|force|version")
	flag.IntVar(&steps, "steps", 1, "number of steps for down/goto")
	flag.Parse()

	migrationsPath := "file://migrations"
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("MYSQL_DSN is empty")
	}

	m, err := migrate.New(migrationsPath, fmt.Sprintf("mysql://%s", dsn))
	if err != nil {
		log.Fatalf("init migrate: %v", err)
	}
	defer func() { _, _ = m.Close() }()

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate up: %v", err)
		}
		log.Println("migrate up done")
	case "down":
		if err := m.Steps(-steps); err != nil {
			log.Fatalf("migrate down %d: %v", steps, err)
		}
		log.Println("migrate down done")
	case "goto":
		if err := m.Steps(steps); err != nil {
			log.Fatalf("migrate steps %d: %v", steps, err)
		}
		log.Println("migrate steps done")
	case "force":
		if err := m.Force(steps); err != nil {
			log.Fatalf("migrate force %d: %v", steps, err)
		}
		log.Println("migrate force done")
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("migrate version: %v", err)
		}
		log.Printf("version: %d dirty=%v\n", v, dirty)
	default:
		log.Fatalf("unknown direction %q", direction)
	}
}
