package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Task    string `json:"task"`
	Content string `json:"content,omitempty"`
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0, // use default DB
	})

	msg := Message{
		From:    "agent_main",
		To:      "agent_planner",
		Task:    "plan_learning",
		Content: "I want to learn AI Agent fundamental",
	}

	jsonMsg, _ := json.Marshal(msg)
	err := rdb.Publish(ctx, "agent-channel", jsonMsg).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("Agent_Main: Message sent!")

mainLoop:
	for {
		// Subscribe to the channel to wait for a reply
		sub := rdb.Subscribe(ctx, "agent-channel")
		ch := sub.Channel()
		fmt.Println("Agent_Main is waiting for a reply...")

		// Wait for a reply from Agent B
		for msg := range ch {
			var message Message
			json.Unmarshal([]byte(msg.Payload), &message)

			if message.To == "agent_main" {
				fmt.Println("Agent_Main received:", message.Task, message.Content)
				break mainLoop // Exit after receiving the reply
			}
		}
	}
}
