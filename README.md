# Telegram Task Manager Bot

A simple Telegram bot written in Go for managing personal tasks directly from Telegram.

## Features

- Add tasks
- List tasks
- Complete tasks
- Delete tasks
- Clear completed tasks
- Task statistics
- Multi-user support
- JSON persistence
- Telegram Bot API

## Tech Stack

- Go
- Telegram Bot API
- JSON

This project uses the official Go Telegram Bot API library to simplify command handling and long-polling updates.

## Bot Commands

| Command | Description |
|---------|-------------|
| /start | Start the bot |
| /help | Show available commands |
| /add <task> | Add a task |
| /list | List your tasks |
| /done <id> | Complete a task |
| /delete <id> | Delete a task |
| /clear | Remove completed tasks |
| /stats | Show task statistics |

## Getting Started

Clone the repository:

```bash
git clone https://github.com/YOUR_USERNAME/telegram-task-bot.git
cd telegram-task-bot
```

## Creating a Telegram Bot

1. Open Telegram.
2. Find @BotFather.
3. Run /newbot.
4. Follow BotFather's instructions.
5. Copy the generated bot token.
6. Never publish the token.

## Configuration

Set the required environment variable:

```bash
export TELEGRAM_BOT_TOKEN="your_token_here"
```

For Windows PowerShell:

```powershell
$env:TELEGRAM_BOT_TOKEN="your_token_here"
```

You can also create a local `.env` file in the project root with the same variable:

```env
TELEGRAM_BOT_TOKEN=your_token_here
```

Make sure `.env` is ignored by Git and never committed.

## Running

```bash
go run .
```

## Running Tests

```bash
go test ./...
```

## Project Structure

```text
telegram-task-bot/
├── main.go
├── internal/
│   ├── bot/
│   │   ├── bot.go
│   │   └── handlers.go
│   ├── storage/
│   │   ├── json.go
│   │   └── json_test.go
│   └── todo/
│       ├── service.go
│       ├── service_test.go
│       └── task.go
├── data/
│   └── .gitkeep
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── data/tasks.json (created at runtime, ignored by Git)
```

The `bot` package handles Telegram commands, the `todo` package contains task logic, and the `storage` package persists tasks to a JSON file.

## Example Usage

User:
```text
/add Learn Go
```

Bot:
```text
✅ Task added!
1. Learn Go
```

User:
```text
/list
```

Bot:
```text
📋 Your tasks

1. ⬜ Learn Go
```

User:
```text
/done 1
```

Bot:
```text
✅ Task completed!

1. Learn Go
```

User:
```text
/stats
```

Bot:
```text
📊 Task Statistics

Total: 1
Completed: 1
Pending: 0
Progress: 100%
```

## What I Learned

- working with Telegram Bot API
- building bots in Go
- processing commands
- managing multiple users
- JSON persistence
- concurrency safety
- environment variables
- error handling
- unit testing

## Future Improvements

- task priorities
- due dates
- reminders
- PostgreSQL storage
- Docker deployment
- inline keyboards
