package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
		DB:   0, // use default DB)
	})

	sub := rdb.Subscribe(ctx, "agent-channel")
	ch := sub.Channel()

	fmt.Println("Agent_Planner is listening for messages...")

	for msg := range ch {
		var message Message
		json.Unmarshal([]byte(msg.Payload), &message)

		if message.To == "agent_planner" {
			fmt.Printf("Agent_Planner received: %s\n", message.Task)

			reply := Message{
				From:    "agent_planner",
				To:      message.From,
				Task:    "plan_learning",
				Content: planLearning(message.Content),
			}

			jsonReply, _ := json.Marshal(reply)
			rdb.Publish(ctx, "agent-channel", jsonReply)
			fmt.Printf("Agent_Planner sent reply: %s\n", reply.Content)
		}
	}
}

func planLearning(goal string) string {
	if strings.Contains(strings.ToLower(goal), "ai") {
		return strings.Join([]string{
			"1. Learn Python basics",
			"2. Understand ML concepts (supervised, unsupervised)",
			"3. Practice with Scikit-learn",
			"4. Try real datasets (e.g., Kaggle)",
			"5. Explore deep learning with PyTorch or TensorFlow",
		}, "\n")
	}

	return "Research topic not recognized for AI learning plan."
}
