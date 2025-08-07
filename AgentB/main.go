package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"os"

	"github.com/joho/godotenv"
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
	// Load .env file
	godotenv.Load()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		panic("Please set GEMINI_API_KEY in .env file")
	}

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
				Content: planLearning(message.Content, apiKey),
			}

			jsonReply, _ := json.Marshal(reply)
			rdb.Publish(ctx, "agent-channel", jsonReply)
			fmt.Printf("Agent_Planner sent reply: %s\n", reply.Content)
		}
	}
}

func callGeminiAPI(prompt string, apiKey string) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + apiKey

	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API error: %v", result)
	}

	if candidates, ok := result["candidates"].([]interface{}); ok && len(candidates) > 0 {
		first := candidates[0].(map[string]interface{})
		content := first["content"].(map[string]interface{})
		parts := content["parts"].([]interface{})
		text := parts[0].(map[string]interface{})["text"].(string)
		return text, nil
	}

	return "Failed to get response from Gemini", nil
}

func planLearning(goal string, apiKey string) string {
	prompt := fmt.Sprintf("Help me create a step-by-step learning plan for: %s", goal)
	response, err := callGeminiAPI(prompt, apiKey)
	if err != nil {
		return "Error occurred while calling Gemini API"
	}

	lines := strings.Split(response, "\n")
	var steps []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			steps = append(steps, line)
		}
	}

	return strings.Join(steps, "\n")
}
