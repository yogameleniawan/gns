package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	containerPkg "github.com/base-go/backend/container"
	"github.com/base-go/backend/pkg/config"
	"github.com/base-go/backend/pkg/server"
)

var environtment string

func init() {
	env := flag.String("env", "development", "Environment (development/production)")
	flag.Parse()

	environtment = *env

	switch *env {
	case "development":
		config.LoadConfig("./config/config.development.yaml")
	case "staging":
		config.LoadConfig("./config/config.staging.yaml")
	case "production":
		config.LoadConfig("./config/config.production.yaml")
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// setup server
	container, err := containerPkg.New()
	if err != nil {
		log.Fatal(err)
	}

	err = container.Invoke(Start)
	if err != nil {
		log.Fatal(err)
	}

	// add quit signal
	quit := make(chan os.Signal, 1)
	container.Provide(func() chan os.Signal {
		return quit
	})

	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	err = container.Invoke(Shutdown)
	if err != nil {
		log.Fatal(err)
	}
}

func Start(
	svr server.Server,
) {

	cfg := config.GetConfig()

	log.Println("Check migration")
	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Error running migration: %s", err)
	}
	log.Println("Database migration succeeded")

	log.Println("Starting server...")
	if err := svr.Start(); err != nil {
		log.Fatal(err)
	}

	if environtment == "production" {
		// for the security purpose, we need to remove file configuration after server start
		// This can be used in case the container where the backend is located is hacked,
		// information related to DB configuration, etc. cannot be accessed.
		log.Println("Deleting configuration after service running...")
		if err := os.Remove("./config/config.production.yaml"); err != nil {
			log.Fatal(err)
		}
	}
}

func Shutdown(
	quit chan os.Signal,
	svr server.Server,
) {

	q := <-quit
	log.Println("got signal:", q)

	log.Println("Shutting down server...")
	if err := svr.Stop(); err != nil {
		log.Fatal(err)
	}

	log.Println("service gracefully shutdown")
}
