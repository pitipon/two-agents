## Running with Docker Compose

You can run all agents and Redis together using Docker Compose:

```sh
docker-compose up --build
```

This will start:
- Redis
- Agent Main (agentA) on port **8081**
- Agent Planner (agentB)
- Agent Critic (agentC)
- Agent Decider (agentD)

Each agent is configured via environment variables in `docker-compose.yml`.

## Sending a Request

Once running, you can interact with the system by sending a POST request to Agent Main:

```sh
curl -X POST http://localhost:8081/run \
  -H "Content-Type: application/json" \
  -d '{"input": "I want to learn React"}'
```

**Response:**
```json
{
  "result": "...finalized plan...",
  "run_id": "123456"
}
```

## How It Works

- Agent Main receives the learning goal via HTTP.
- Agent Planner generates a learning plan.
- Agent Critic reviews and suggests improvements.
- Agent Decider finalizes the plan and returns it to Agent Main.

---

Feel free to modify the agents or the workflow to suit your needs!