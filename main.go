package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Agent interface {
	Run(ctx context.Context, input string) (string, error)
}

type LlmAgent struct {
	Name string
	Role string
}

func (a *LlmAgent) Run(ctx context.Context, input string) (string, error) {
	// TODO: call an llm here (gemini, openai, etc.)
	out := callLLM(ctx, a.Name, a.Role, input)
	return out, nil
}

func callLLM(ctx context.Context, name, role, input string) string {
	time.Sleep(time.Duration(100 * time.Millisecond)) // Simulate LLM processing time
	return "LLM response for " + name + " (" + role + "): " + input
}

// Redis async helpers
func publishToChannel(ctx context.Context, rdb *redis.Client, ch string, msg string) error {
	return rdb.Publish(ctx, ch, msg).Err()
}

func extractRunID(channel string) string {
	parts := strings.Split(channel, ":")
	if len(parts) >= 4 {
		return parts[3]
	}
	return "unknown"
}

func startSingleAgentWorker(ctx context.Context, rdb *redis.Client, agentID, nextID string, ag Agent) {
	pattern := fmt.Sprintf("agent:%s:in:*", agentID)
	sub := rdb.PSubscribe(ctx, pattern)
	defer sub.Close()

	log.Println("Agent %s Listening on pattern: %s", agentID, pattern)
	log.Println("nextID %s", nextID)

	for msg := range sub.Channel() {
		fmt.Println(msg)
		runID := extractRunID(msg.Channel)
		log.Printf("[%s] received input: %s", agentID, msg.Payload)
		out, err := ag.Run(ctx, msg.Payload)
		if err != nil {
			log.Printf("[%s] error: %v", agentID, err)
			continue
		}
		if nextID == "" {
			finalCh := fmt.Sprintf("flow:result:%s", runID)
			_ = publishToChannel(ctx, rdb, finalCh, out)
			log.Println("[%s] published final result to %s", agentID, finalCh)
		} else {
			nextCh := fmt.Sprintf("agent:%s:in:%s", nextID, runID)
			_ = publishToChannel(ctx, rdb, nextCh, out)
			log.Printf("[%s] published output to next agent channel %s", agentID, nextCh)
		}

	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	agentID := os.Getenv("AGENT_ID")
	if agentID == "" {
		log.Fatal("AGENT_ID env var is required")
	}

	agentRole := os.Getenv("AGENT_ROLE")
	if agentRole == "" {
		log.Fatal("AGENT_ROLE env var is required")
	}

	// it can be empty for last agent
	nextAgentID := os.Getenv("NEXT_AGENT_ID")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	agent := &LlmAgent{
		Name: agentID,
		Role: agentRole,
	}

	// Start worker
	go startSingleAgentWorker(ctx, rdb, agentID, nextAgentID, agent)

	// HTTP server - single endpoint to trigger flow by publishging to first agent input channel
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/run", func(c *gin.Context) {
		var body struct {
			Input string `json:"input"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if agentID != "agentA" {
			c.JSON(400, gin.H{"error": "Only agentA container can receive initial request"})
			return
		}

		runID := fmt.Sprint(rand.Intn(1000000))
		firstCh := fmt.Sprintf("agent:%s:in:%s", agentID, runID)
		fmt.Println(firstCh)
		fmt.Println(body.Input)

		if err := publishToChannel(ctx, rdb, firstCh, body.Input); err != nil {
			c.JSON(500, gin.H{"error": "Failed to publish to Redis"})
			return
		}

		finalCh := fmt.Sprintf("flow:result:%s", runID)
		sub := rdb.Subscribe(ctx, finalCh)
		defer sub.Close()
		ch := sub.Channel()

		select {
		case msg := <-ch:
			c.JSON(200, gin.H{"result": msg.Payload, "run_id": runID})
		case <-time.After(15 * time.Second):
			c.JSON(504, gin.H{"error": "Timeout waiting for async result"})
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Agent %s starting HTTP server on: %s", agentID, port)
	r.Run(":" + port)
}
