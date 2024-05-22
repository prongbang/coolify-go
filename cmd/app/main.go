package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

func main() {
	app := fiber.New()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0, // use default DB
	})

	ctx := context.Background()
	err := rdb.Set(ctx, "key", "Hello, Coolify with GO!", 0).Err()
	if err != nil {
		log.Error(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		val, err := rdb.Get(ctx, "key").Result()
		if err != nil {
			log.Error(err)
			return c.SendString(err.Error())
		}
		return c.SendString(val)
	})

	log.Fatal(app.Listen(":9001"))
}
