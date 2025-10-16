# Subscriptions Service

Subscriptions management service

## Project structure

- `cmd/` — where main.go is located
- `internal/` — private packages
- `pkg/database/` — connection PostgreSQL
- `Dockerfile` — docker container
- `docker-compose.yml` — run all neccessary containers together
- `Makefile` — simplified init, test, build commands
- `.env` — environment

## Run project

### 1. Clone repo

```bash
git clone https://github.com/teamcutter/subscriptions-service-task.git
cd subscriptions-service-task
```

### 2. Create .env
Idealy you have to create your own .env file, but I have pushed them for you JUST for simplicity purposes. It is bad to push .env files, I know :D

### 3. Reminder
Do not forget to run Docker on your machine and install Go 

### 4. Start project
Simply do
```bash
make up-build
```