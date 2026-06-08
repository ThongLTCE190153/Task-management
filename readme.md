# Task Management API - Golang

## Overview

This project is a simple Task Management REST API built with Golang using Clean Architecture principles.

The project was developed for learning backend development concepts including:

* REST API
* Clean Architecture
* PostgreSQL
* Docker
* JWT Authentication
* Repository Pattern
* Middleware
* Password Hashing
* Migration

---

# Technologies

* Golang
* Gin Gonic
* PostgreSQL
* Docker
* JWT
* bcrypt
* golang-migrate

---

# Project Structure

```bash
task-golang/
│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── config/
│   ├── database/
│   ├── dto/
│   ├── entities/
│   ├── handlers/
│   ├── mappers/
│   ├── middlewares/
│   ├── repositories/
│   ├── responses/
│   ├── routes/
│   ├── services/
│   └── utils/
│
├── migrations/
│
├── .env
├── docker-compose.yml
├── go.mod
└── README.md
```

---

# Features

## Authentication

* Register user
* Login user
* Password hashing using bcrypt
* JWT token authentication

## User

* Create user
* Get user by ID

## Project

* Create project
* Get all projects
* Get project by ID

## Task

* Create task
* Get all tasks
* Get task by ID

---

# Database Setup

## Run PostgreSQL using Docker

```bash
docker-compose up -d
```

Check running containers:

```bash
docker ps
```

---

# Environment Variables

Create `.env` file:

```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=task_management

JWT_SECRET=your_long_random_secret
```

---

# Run Migration

Install migrate:

```bash
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Run migration:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5433/task_management?sslmode=disable" up
```

---

# Install Dependencies

```bash
go mod tidy
```

---

# Run Project

```bash
go run ./cmd
```

Server runs at:

```bash
http://localhost:8080
```

---

# API Endpoints

# Authentication

## Register

### POST /auth/register

Request:

```json
{
  "name": "Thong",
  "email": "thong@gmail.com",
  "password": "123456"
}
```

Response:

```json
{
  "success": true,
  "message": "Register success",
  "data": null
}
```

---

## Login

### POST /auth/login

Request:

```json
{
  "email": "thong@gmail.com",
  "password": "123456"
}
```

Response:

```json
{
  "success": true,
  "message": "Login success",
  "data": {
    "token": "jwt_token_here"
  }
}
```

---

# Authorization

Protected APIs require JWT token.

Add header:

```http
Authorization: Bearer your_token
```

---

# Project APIs

## Create Project

### POST /projects

Request:

```json
{
  "name": "Task Management",
  "description": "Backend API Project"
}
```

---

## Get All Projects

### GET /projects

---

## Get Project By ID

### GET /projects/:id

---

# Task APIs

## Create Task

### POST /tasks

Request:

```json
{
  "title": "Learn Golang",
  "description": "Study Clean Architecture",
  "status": "todo",
  "project_id": 1
}
```

---

## Get All Tasks

### GET /tasks

---

## Get Task By ID

### GET /tasks/:id

---

# Middleware

This project uses middleware for:

* JWT authentication
* Request validation
* Error handling

---

# Clean Architecture

This project separates code into multiple layers:

| Layer      | Responsibility                      |
| ---------- | ----------------------------------- |
| Handler    | Handle HTTP request/response        |
| Service    | Business logic                      |
| Repository | Database logic                      |
| Entity     | Database model                      |
| DTO        | Request/response object             |
| Middleware | Authentication & request processing |

---

# Future Improvements

* Update/Delete APIs
* Role-based authorization
* Unit testing
* Swagger documentation
* Refresh token
* Pagination
* Search & filtering

---

# Author

Thong
