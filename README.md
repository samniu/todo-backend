# Todo App Backend

A Go-powered backend service for a Todo application with RESTful API and real-time WebSocket notifications.

*[中文文档](README_ZH.md)*

## Overview

This backend service provides a robust API for a Todo application, featuring user authentication, task management, and real-time notifications. It's designed to support both online and offline operations, allowing for seamless synchronization when network connection is restored.

## Features

- 🔒 **User Authentication**: Secure registration and login system with JWT
- ✅ **Task Management**: Complete CRUD operations for todos
- 🔔 **Real-time Notifications**: WebSocket support for instant updates
- 🔄 **Offline Sync Support**: Special API design to handle client offline operations

## Getting Started

### Prerequisites

- Go 1.18+
- PostgreSQL/MySQL (configurable)

### Installation

1. Clone the repository

```bash
git clone https://github.com/samniu/todo-backend.git
cd todo-backend
```

2. Set up environment variables

Create a `.env` file:

```
DB_CONNECTION=postgres://username:password@localhost:5432/todo_db
JWT_SECRET=your_jwt_secret_key
PORT=8080
```

3. Run the server

```bash
go run cmd/main.go
```

The server will start at `http://localhost:8080`.

## License

This project is licensed under the [MIT License](LICENSE).
