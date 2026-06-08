# Task Management System

Hệ thống quản lý công việc nội bộ xây dựng bằng Go, theo kiến trúc Clean Architecture.

## Công nghệ sử dụng

- **Go 1.26** — Ngôn ngữ chính
- **Gin** — Web framework
- **PostgreSQL** — Database
- **Redis** — Caching & Job Queue
- **WebSocket** — Realtime
- **JWT** — Authentication
- **Docker** — Containerization

## Tính năng

- Đăng ký, đăng nhập, JWT authentication, Refresh Token
- Phân quyền theo role: Admin, Manager, User
- Quản lý Project và Task
- Giao task cho user, notification qua background worker
- Comment trong task, realtime qua WebSocket
- Redis caching, Rate limiting
- Pagination & Filtering
- Unit test, Integration test, CI/CD

## Cấu trúc thư mục
```text
Task/
│
├── cmd/
│   ├── api/
│   │   └── main.go
│   │
│   └── worker/
│       └── main.go
│
├── internal/
│   │
│   ├── cache/
│   │   └── redis.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── db/
│   │   └── postgres.go
│   │
│   ├── dto/
│   │   ├── auth_dto.go
│   │   ├── user_dto.go
│   │   ├── project_dto.go
│   │   └── task_dto.go
│   │
│   ├── entities/
│   │   ├── user.go
│   │   ├── project.go
│   │   ├── task.go
│   │   └── notification.go
│   │
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── project_handler.go
│   │   ├── task_handler.go
│   │   └── websocket_handler.go
│   │
│   ├── jobs/
│   │   ├── notification_job.go
│   │   └── task_job.go
│   │
│   ├── mappers/
│   │   ├── user_mapper.go
│   │   ├── project_mapper.go
│   │   └── task_mapper.go
│   │
│   ├── middlewares/
│   │   ├── jwt_middleware.go
│   │   ├── role_middleware.go
│   │   ├── ratelimit_middleware.go
│   │   ├── logger_middleware.go
│   │   └── requestid_middleware.go
│   │
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── project_repository.go
│   │   ├── task_repository.go
│   │   └── notification_repository.go
│   │
│   ├── responses/
│   │   └── response.go
│   │
│   ├── routes/
│   │   ├── auth_routes.go
│   │   ├── user_routes.go
│   │   ├── project_routes.go
│   │   ├── task_routes.go
│   │   └── websocket_routes.go
│   │
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── project_service.go
│   │   ├── task_service.go
│   │   └── notification_service.go
│   │
│   ├── tests/
│   │   ├── auth_service_test.go
│   │   ├── user_service_test.go
│   │   ├── project_service_test.go
│   │   └── task_service_test.go
│   │
│   ├── utils/
│   │   ├── jwt.go
│   │   ├── password.go
│   │   ├── validator.go
│   │   └── time.go
│   │
│   └── websocket/
│       ├── hub.go
│       ├── client.go
│       └── message.go
│
├── migrations/
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_projects.up.sql
│   ├── 000002_create_projects.down.sql
│   ├── 000003_create_tasks.up.sql
│   ├── 000003_create_tasks.down.sql
│   ├── 000004_create_notifications.up.sql
│   └── 000004_create_notifications.down.sql
│
├── .github/
│   └── workflows/
│       └── ci.yml
│
├── docs/
│   └── postman_collection.json
│
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── .env
├── .env.example
├── .gitignore
└── README.md
```

## Yêu cầu

- Docker & Docker Compose

## Cài đặt và chạy

### 1. Clone repository

```bash
git clone https://github.com/ThongLTCE190153/Task-management.git
cd Task-management
```

### 2. Cấu hình môi trường

```bash
cp .env.example .env
```

Mở file `.env` và điền thông tin:

```env
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=task_management
DB_SSLMODE=disable
JWT_SECRET=your_secret_key
JWT_EXPIRE_HOURS=24
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 3. Chạy với Docker Compose

```bash
docker-compose up --build
```

Hệ thống sẽ khởi động:
- API server: http://localhost:8080
- PostgreSQL: localhost:5433
- Redis: localhost:6379

### 4. Chạy migration

```bash
# Chạy từng file migration theo thứ tự
Get-Content migrations/001_create-users.up.sql | docker exec -i task-management-postgres psql -U postgres -d task_management
Get-Content migrations/002_create-projects.up.sql | docker exec -i task-management-postgres psql -U postgres -d task_management
Get-Content migrations/003_create-tasks.up.sql | docker exec -i task-management-postgres psql -U postgres -d task_management
Get-Content migrations/004_create-task-comments.up.sql | docker exec -i task-management-postgres psql -U postgres -d task_management
Get-Content migrations/005_create_role_to_users.up.sql | docker exec -i task-management-postgres psql -U postgres -d task_management
Get-Content migrations/006_create-account_admin.up.sql | docker exec -i task-management-postgres psql -U postgres -d task_management
```

### 5. Tài khoản Admin mặc định
Email:    admin@system.com
Password: admin123456

## API Endpoints

> **Chú thích:**
> - ✅ — Cần JWT token trong header: `Authorization: Bearer <token>`
> - ❌ — Không cần token (public API)
> - **Role:** Admin > Manager > User

### Auth
| Method | Endpoint | Mô tả | Auth |
|--------|----------|-------|------|
| POST | /auth/register | Đăng ký | ❌ |
| POST | /auth/login | Đăng nhập | ❌ |
| POST | /auth/refresh | Refresh token | ❌ |
| POST | /auth/logout | Đăng xuất | ✅ |
| GET | /auth/me | Thông tin user | ✅ |

### Users
| Method | Endpoint | Mô tả | Role |
|--------|----------|-------|------|
| GET | /users | Danh sách user | Admin |
| PUT | /users/:id/role | Đổi role | Admin |
| GET | /users/:id | Chi tiết user | ✅ |
| PUT | /users/:id | Cập nhật user | ✅ |

### Projects
| Method | Endpoint | Mô tả | Role |
|--------|----------|-------|------|
| POST | /projects | Tạo project | Manager+ |
| GET | /projects | Danh sách project | ✅ |
| GET | /projects/:id | Chi tiết project | ✅ |
| PUT | /projects/:id | Cập nhật project | Manager+ |
| DELETE | /projects/:id | Xóa project | Manager+ |

### Tasks
| Method | Endpoint | Mô tả | Role |
|--------|----------|-------|------|
| POST | /tasks | Tạo task | Manager+ |
| GET | /tasks | Danh sách task | ✅ |
| GET | /tasks/:id | Chi tiết task | ✅ |
| PUT | /tasks/:id | Cập nhật task | Manager+ |
| DELETE | /tasks/:id | Xóa task | ✅ |

### Comments
| Method | Endpoint | Mô tả | Auth |
|--------|----------|-------|------|
| POST | /comments/task/:taskId | Tạo comment | ✅ |
| GET | /comments/task/:taskId | Danh sách comment | ✅ |

### WebSocket
| Endpoint | Mô tả |
|----------|-------|
| GET /ws?token=xxx | Kết nối realtime |

### Health Check
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | /health | Kiểm tra trạng thái |

## Chạy test

```bash
# Unit tests
go test ./internal/services/... -v

# Integration tests (cần PostgreSQL và Redis đang chạy)
go test ./internal/tests/... -v

# Race detector
go test -race ./internal/services/... -v
```

## WebSocket Events

| Event | Trigger |
|-------|---------|
| task.created | Khi tạo task mới |
| task.status_updated | Khi đổi trạng thái task |
| task.deleted | Khi xóa task |
| task.comment.created | Khi có comment mới |
| project.created | Khi tạo project |
| project.updated | Khi cập nhật project |
| project.deleted | Khi xóa project |
| task.assigned | Khi assign task (qua worker) |