package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"

	"github.com/Dolan-dev/URL-Shortener/config"
	"github.com/Dolan-dev/URL-Shortener/handler"
	"github.com/Dolan-dev/URL-Shortener/url"
)

func main() {
	cfg := &config.Config{}
	err := envconfig.Process("shorturl", cfg)
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	shortener := url.NewShortener(redisClient, cfg)

	create := handler.NewCreate(shortener, cfg)
	redirect := handler.NewRedirect(shortener)

	r := chi.NewRouter()

	r.Get("/", handler.Index)
	r.Post("/", create.Handler)
	r.Get("/{shortURL}", redirect.Handler)

	http.ListenAndServe(":8080", r)
}
