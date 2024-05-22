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

	//redisURL := "redis://user:password@localhost:6379/0?protocol=3"
	redisURL := os.Getenv("REDIS_URL")

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Error(err)
	}
	rdb := redis.NewClient(opts)

	ctx := context.Background()
	err = rdb.Set(ctx, "key", "Hello, Coolify with GO!", 0).Err()
	if err != nil {
		log.Error(err)
	}

	app.Get("/reload", func(c *fiber.Ctx) error {
		redisURL = os.Getenv("REDIS_URL")

		opts, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Error(err)
			return c.SendString(err.Error())
		}
		rdb = redis.NewClient(opts)

		return c.SendString("OK")
	})

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"url": redisURL,
		})
	})

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
