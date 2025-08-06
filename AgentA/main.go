package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Message struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0, // use default DB
	})

	msg := Message{
		From: "AgentA",
		To:   "AgentB",
		Text: "Hello from Agent A!",
	}

	jsonMsg, _ := json.Marshal(msg)
	err := rdb.Publish(ctx, "agent-channel", jsonMsg).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Agent A: Message sent!")
}
