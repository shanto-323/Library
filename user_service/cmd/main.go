package main

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/shanto-323/Library/books"
	userservice "github.com/shanto-323/Library/user_service"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseDsn string `envconfig:"DATABASE_DSN"`
}

func main() {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		books.LogError(slog.LevelError, "MAIN", err, fmt.Sprintf("failed to load config %s", cfg.DatabaseDsn))
		return
	}

	var r userservice.Repository
	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			r, err = userservice.NewUserRepostiry(cfg.DatabaseDsn)
			if err != nil {
				books.LogError(slog.LevelError, "MAIN", err, fmt.Sprintf("failed to get reposiroty with dsn %s", cfg.DatabaseDsn))
				return err
			}
			return nil
		},
	)

	books.LogInfo(slog.LevelInfo, "MAIN", "User-Service running on port 8080 ...")
	s := userservice.NewUserService(r)
	log.Fatal(userservice.ListenGRPC(s, ":8080"))
}
