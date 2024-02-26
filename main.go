package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
	redisHost := os.Getenv("REDIS_HOST")
	if len(redisHost) == 0 {
		panic("no REDIS_HOST env")
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}
}
func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/counter", func(c *fiber.Ctx) error {
		var rdKey = "ct"
		var count int
		var err error
		countInDB, err := rdb.Get(ctx, rdKey).Result()
		if err == redis.Nil {
			fmt.Println("key \"" + rdKey + "\" does not exist")
		} else if err != nil {
			fmt.Println(err.Error())
		} else {
			count, err = strconv.Atoi(countInDB)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		newCount := count + 1
		err = rdb.Set(ctx, rdKey, strconv.Itoa(newCount), 0).Err()
		if err != nil {
			fmt.Println(err.Error())
		}

		return c.SendString("Count: " + strconv.Itoa(newCount))
	})

	app.Listen(":3000")
}
