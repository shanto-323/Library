package main

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/shanto-323/Library/books"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseDsn string `envconfig:"DATABASE_DSN"`
}

func main() {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r books.Repository
	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			r, err = books.NewRepository(cfg.DatabaseDsn)
			if err != nil {
				fmt.Println(err)
				return err
			}
			return nil
		},
	)

	books.LogInfo(slog.LevelInfo, "MAIN", "User-Service running on port 8080 ...")
	s := books.NewBookService(r)
	log.Fatal(books.ListenGRPC(s, ":8080"))
}
