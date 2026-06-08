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
Task/
├── cmd/
│   ├── api/main.go        # API server
│   └── worker/main.go     # Notification worker
├── internal/
│   ├── cache/             # Redis connection
│   ├── config/            # App configuration
│   ├── db/                # Database connection
│   ├── dto/               # Request/Response objects
│   ├── entities/          # Domain models
│   ├── handlers/          # HTTP handlers
│   ├── jobs/              # Background jobs
│   ├── mappers/           # Entity <-> DTO mappers
│   ├── middlewares/       # Auth, Role, RateLimit, Logger
│   ├── repositories/      # Database layer
│   ├── responses/         # Response format
│   ├── routes/            # API routes
│   ├── services/          # Business logic
│   ├── tests/             # Integration tests
│   ├── utils/             # JWT, Password utils
│   └── websocket/         # WebSocket Hub
├── migrations/            # Database migrations
├── .github/workflows/     # CI/CD pipeline
├── Dockerfile
├── docker-compose.yml
└── .env.example
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