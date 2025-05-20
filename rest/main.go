package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BookClient        string `envconfig:"BOOK_SERVICE_URL"`
	UserServiceClient string `envconfig:"USER_SERVICE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewServer(":8080", cfg.BookClient, cfg.UserServiceClient)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.Start())
}
