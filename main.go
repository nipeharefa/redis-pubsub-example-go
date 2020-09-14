package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func consume(pubsub *redis.PubSub) {
	var err error
	// Wait for confirmation that subscription is created before publishing anything.
	_, err = pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()
	// Consume messages.
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
	}
}

func main() {

	var err error

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()
	fmt.Println(pong, err)

	pubsub := rdb.Subscribe(ctx, "mychannel1")

	args := os.Args

	// fmt.Println(args[1])
	if len(args) >= 2 {
		switch args[1] {
		case "sub":
			fmt.Println("Sub")
			consume(pubsub)
			return
		}
	}

	fmt.Println("Pub")
	// Publish a message.
	err = rdb.Publish(ctx, "mychannel1", "hello").Err()
	if err != nil {
		panic(err)
	}
}
