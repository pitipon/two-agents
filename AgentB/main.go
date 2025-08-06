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
		DB:   0, // use default DB)
	})

	sub := rdb.Subscribe(ctx, "agent-channel")
	ch := sub.Channel()

	fmt.Println("AgentB is listening for messages...")

	for msg := range ch {
		var message Message
		json.Unmarshal([]byte(msg.Payload), &message)

		if message.To == "AgentB" {
			fmt.Printf("Agent B received: %s\n", message.Text)

			reply := Message{
				From: "AgentB",
				To:   "AgentA",
				Text: "Hello Agent A, I got your message!",
			}

			jsonReply, _ := json.Marshal(reply)
			rdb.Publish(ctx, "agent-channel", jsonReply)
		}
	}
}
