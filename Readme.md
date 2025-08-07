# Two-Agents System

This project demonstrates a multi-agent system using Redis Pub/Sub and Google Gemini API. There are **four agents** that communicate in a pipeline:

1. **Agent Main**: Starts the process by sending a learning goal.
2. **Agent Planner**: Creates a step-by-step learning plan.
3. **Agent Critic**: Critiques or suggests improvements to the plan.
4. **Agent Decider**: Finalizes the plan and sends it back to the main agent.

## Architecture

```
Agent Main → Agent Planner → Agent Critic → Agent Decider → Agent Main
```

## Prerequisites

- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- A Google Gemini API key (set as `GEMINI_API_KEY` in a `.env` file in each agent's folder)

## Run Redis

Start Redis using Docker:

```
docker run -p 6379:6379 redis
```

## Running the Agents

Open **four terminals**, one for each agent:

**Terminal 1: Agent Main**
```
cd AgentA
go run main.go
```

**Terminal 2: Agent Planner**
```
cd AgentB
go run main.go
```

**Terminal 3: Agent Critic**
```
cd AgentC
go run main.go
```

**Terminal 4: Agent Decider**
```
cd AgentD
go run main.go
```

## How It Works

- Agent Main sends a learning goal (e.g., "I want to learn React").
- Agent Planner generates a learning plan.
- Agent Critic reviews and suggests improvements.
- Agent Decider finalizes the plan and returns it to Agent Main.

---

Feel free to modify the agents or the workflow to suit your needs!