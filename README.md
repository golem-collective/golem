# Golem Agent

A Go-based AI agent application that provides a conversational interface powered by OpenAI's API. The application stores conversation history in PostgreSQL, making it sessionless and persistent.

## Features

- Console-based chat interface
- PostgreSQL database integration for persistent storage
- OpenAI API integration for natural language processing
- Session-independent conversation history
- Simple and clean architecture

## Prerequisites

- Go 1.18 or higher
- PostgreSQL 12 or higher
- OpenAI API key

## Setup

### 1. Clone the repository

```bash
git clone https://github.com/golem-collective/golem.git
```

### 2. Set up the database

Create a PostgreSQL database for the application:

```bash
psql -U postgres
CREATE DATABASE go;
\q
```

### 3. Configure environment variables

Create a `.env` file in the project root with the following variables:

```
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=go
DB_HOST=localhost
DB_PORT=5432
OPENAI_API_KEY=your_openai_api_key
```

Replace `your_password` with your PostgreSQL password and `your_openai_api_key` with your OpenAI API key.

### 4. Install dependencies

```bash
go mod tidy
```

### 5. Build the application

```bash
go build
```

## Usage

### Console Interface

Run the application:

```bash
./ai-agent-app
```

This will start the console interface where you can chat with the AI agent. The application will:

1. Connect to the PostgreSQL database
2. Create necessary tables if they don't exist
3. Create a default agent
4. Start the chat interface

Commands:
- Type your message and press Enter to chat with the agent
- Type `clear` to clear the conversation history
- Type `exit` to quit the application

### API Integration

The application also provides HTTP endpoints for integration with other applications:

- `POST /api/agents` - Create a new agent
- `POST /api/agents/{agentID}/chat` - Chat with an agent

## Project Structure

- `main.go` - Application entry point
- `models/` - Data models
- `handlers/` - HTTP request handlers
- `services/` - Business logic
- `database/` - Database connection and operations

## Database Schema

### Agents Table

Stores information about AI agents:

```sql
CREATE TABLE agents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    context TEXT
);
```

### Chat History Table

Stores conversation history:

```sql
CREATE TABLE chat_history (
    id SERIAL PRIMARY KEY,
    agent_id INTEGER NOT NULL,
    role VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (agent_id) REFERENCES agents(id)
);
```

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 
