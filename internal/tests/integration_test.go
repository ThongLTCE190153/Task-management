package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"trithong.com/task-golang/internal/cache"
	"trithong.com/task-golang/internal/config"
	"trithong.com/task-golang/internal/db"
	"trithong.com/task-golang/internal/handlers"
	"trithong.com/task-golang/internal/repositories"
	"trithong.com/task-golang/internal/routes"
	"trithong.com/task-golang/internal/services"
	"trithong.com/task-golang/internal/utils"
	appwebsocket "trithong.com/task-golang/internal/websocket"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ===== SETUP SUITE =====
// Giống @SpringBootTest trong Spring Boot

type IntegrationTestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *sql.DB
	token  string
}

func (s *IntegrationTestSuite) SetupSuite() {
	// Load config
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "123456")
	os.Setenv("DB_NAME", "task_management")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("JWT_SECRET", "test-secret")

	cfg := config.LoadConfig()

	utils.SetJWTConfig(cfg.JWTSecret, cfg.JWTExpireHours)

	// Kết nối database thật
	database := db.ConnectPostgres(cfg)
	redisClient := cache.ConnectRedis(cfg)
	hub := appwebsocket.NewHub()

	// Setup repositories
	userRepo := repositories.NewUserRepository(database)
	projectRepo := repositories.NewProjectRepository(database)
	taskRepo := repositories.NewTaskRepository(database)
	commentRepo := repositories.NewCommentRepository(database)

	// Setup services
	authService := services.NewAuthService(userRepo, redisClient)
	projectService := services.NewProjectService(projectRepo, redisClient, hub)
	taskService := services.NewTaskService(taskRepo, redisClient, hub)
	commentService := services.NewCommentService(commentRepo, hub)
	userService := services.NewUserService(userRepo)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(authService)
	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService)
	commentHandler := handlers.NewCommentHandler(commentService)
	userHandler := handlers.NewUserHandler(userService)
	webSocketHandler := handlers.NewWebSocketHandler(hub)

	// Setup router
	gin.SetMode(gin.TestMode)
	s.router = gin.New()

	routes.AuthRoutes(s.router, authHandler)
	routes.UserRoutes(s.router, userHandler)
	routes.ProjectRoutes(s.router, projectHandler, redisClient)
	routes.TaskRoutes(s.router, taskHandler, redisClient)
	routes.CommentRoutes(s.router, commentHandler, redisClient)
	routes.WebSocketRoutes(s.router, webSocketHandler)

	s.db = database

	s.db.Exec(`
		UPDATE users 
		SET role = 'manager' 
		WHERE email = 'test_user@gmail.com'
	`)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	// Xóa data test sau khi chạy xong
	s.db.Exec("DELETE FROM task_comments WHERE content LIKE 'Test%'")
	s.db.Exec("DELETE FROM tasks WHERE title LIKE 'Test%'")
	s.db.Exec("DELETE FROM projects WHERE name LIKE 'Test%'")
	s.db.Exec("DELETE FROM users WHERE email LIKE 'test_%@gmail.com'")
}

// ===== HELPER =====

func (s *IntegrationTestSuite) makeRequest(
	method string,
	url string,
	body interface{},
	token string,
) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	return w
}

// ===== TEST FLOW =====

func (s *IntegrationTestSuite) Test01_Register() {
	w := s.makeRequest("POST", "/auth/register", map[string]string{
		"email":     "test_user@gmail.com",
		"password":  "123456",
		"full_name": "Test User",
	}, "")

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(s.T(), response["data"])
}

func (s *IntegrationTestSuite) Test02_Login() {
	w := s.makeRequest("POST", "/auth/login", map[string]string{
		"email":    "test_user@gmail.com",
		"password": "123456",
	}, "")

	assert.Equal(s.T(), http.StatusOK, w.Code)
	s.db.Exec(
		"UPDATE users SET role = 'manager' WHERE email = 'test_user@gmail.com'",
	)

	w = s.makeRequest("POST", "/auth/login", map[string]string{
		"email":    "test_user@gmail.com",
		"password": "123456",
	}, "")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	s.token = data["token"].(string)

	assert.NotEmpty(s.T(), s.token)
}

func (s *IntegrationTestSuite) Test03_CreateProject() {
	w := s.makeRequest("POST", "/projects", map[string]string{
		"name":        "Test Project",
		"description": "Test Description",
	}, s.token)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(s.T(), "Test Project", data["name"])
}

func (s *IntegrationTestSuite) Test04_CreateTask() {
	// Lấy project id trước
	var projectID int
	s.db.QueryRow(
		"SELECT id FROM projects WHERE name = 'Test Project' LIMIT 1",
	).Scan(&projectID)

	w := s.makeRequest("POST", "/tasks", map[string]interface{}{
		"project_id":  projectID,
		"title":       "Test Task",
		"description": "Test Description",
		"status":      "todo",
	}, s.token)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(s.T(), "Test Task", data["title"])
	assert.Equal(s.T(), "todo", data["status"])
}

func (s *IntegrationTestSuite) Test05_AddComment() {
	// Lấy task id trước
	var taskID int
	s.db.QueryRow(
		"SELECT id FROM tasks WHERE title = 'Test Task' LIMIT 1",
	).Scan(&taskID)

	w := s.makeRequest(
		"POST",
		fmt.Sprintf("/comments/task/%d", taskID),
		map[string]string{
			"content": "Test Comment",
		},
		s.token,
	)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(s.T(), "Test Comment", data["content"])
}

func (s *IntegrationTestSuite) Test06_UpdateTaskStatus() {
	// Lấy task id trước
	var taskID int
	s.db.QueryRow(
		"SELECT id FROM tasks WHERE title = 'Test Task' LIMIT 1",
	).Scan(&taskID)

	w := s.makeRequest(
		"PUT",
		fmt.Sprintf("/tasks/%d", taskID),
		map[string]string{
			"title":  "Test Task",
			"status": "in_progress",
		},
		s.token,
	)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(s.T(), "in_progress", data["status"])
}

func (s *IntegrationTestSuite) Test07_GetTaskList() {
	w := s.makeRequest("GET", "/tasks", nil, s.token)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotNil(s.T(), response["data"])
}

// ===== RUN SUITE =====

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// ===== W4-T03: UNAUTHORIZED TESTS =====

func (s *IntegrationTestSuite) Test08_Unauthorized_NoToken() {
	// Gọi API không có token → phải trả về 401
	w := s.makeRequest("GET", "/tasks", nil, "")

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
}

func (s *IntegrationTestSuite) Test09_Unauthorized_WrongToken() {
	// Gọi API với token sai → phải trả về 401
	w := s.makeRequest("GET", "/tasks", nil, "wrong-token")

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
}

func (s *IntegrationTestSuite) Test10_Unauthorized_NoToken_CreateProject() {
	// Tạo project không có token → phải trả về 401
	w := s.makeRequest("POST", "/projects", map[string]string{
		"name": "Test Project",
	}, "")

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
}

// ===== W4-T04: PERMISSION TESTS =====

func (s *IntegrationTestSuite) Test11_Permission_UpdateOtherUserTask() {
	// Tạo user thứ 2
	s.makeRequest("POST", "/auth/register", map[string]string{
		"email":     "test_user2@gmail.com",
		"password":  "123456",
		"full_name": "Test User 2",
	}, "")

	// Login user 2 lấy token
	w := s.makeRequest("POST", "/auth/login", map[string]string{
		"email":    "test_user2@gmail.com",
		"password": "123456",
	}, "")

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	data := loginResponse["data"].(map[string]interface{})
	token2 := data["token"].(string)

	// Lấy task của user 1
	var taskID int
	s.db.QueryRow(
		"SELECT id FROM tasks WHERE title = 'Test Task' LIMIT 1",
	).Scan(&taskID)

	// User 2 thử update task của user 1 → phải trả về 404
	w = s.makeRequest(
		"PUT",
		fmt.Sprintf("/tasks/%d", taskID),
		map[string]string{
			"title":  "Hacked Task",
			"status": "done",
		},
		token2,
	)

	assert.Equal(s.T(), http.StatusForbidden, w.Code)

}

func (s *IntegrationTestSuite) Test12_Permission_DeleteOtherUserProject() {
	// Lấy project của user 1
	var projectID int
	s.db.QueryRow(
		"SELECT id FROM projects WHERE name = 'Test Project' LIMIT 1",
	).Scan(&projectID)

	// Login user 2
	w := s.makeRequest("POST", "/auth/login", map[string]string{
		"email":    "test_user2@gmail.com",
		"password": "123456",
	}, "")

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	data := loginResponse["data"].(map[string]interface{})
	token2 := data["token"].(string)

	// User 2 thử xóa project của user 1 → phải trả về 404
	w = s.makeRequest(
		"DELETE",
		fmt.Sprintf("/projects/%d", projectID),
		nil,
		token2,
	)

	assert.Equal(s.T(), http.StatusForbidden, w.Code)
}

func (s *IntegrationTestSuite) Test13_RefreshToken() {
	// Login lấy refresh token
	w := s.makeRequest("POST", "/auth/login", map[string]string{
		"email":    "test_user@gmail.com",
		"password": "123456",
	}, "")

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	data := loginResponse["data"].(map[string]interface{})
	refreshToken := data["refresh_token"].(string)

	// Dùng refresh token lấy access token mới
	w = s.makeRequest("POST", "/auth/refresh", map[string]string{
		"refresh_token": refreshToken,
	}, "")

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var refreshResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &refreshResponse)
	refreshData := refreshResponse["data"].(map[string]interface{})

	assert.NotEmpty(s.T(), refreshData["token"])
	assert.NotEmpty(s.T(), refreshData["refresh_token"])
}

func (s *IntegrationTestSuite) Test14_Logout() {
	// Login lấy token
	w := s.makeRequest("POST", "/auth/login", map[string]string{
		"email":    "test_user@gmail.com",
		"password": "123456",
	}, "")

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	data := loginResponse["data"].(map[string]interface{})
	token := data["token"].(string)
	refreshToken := data["refresh_token"].(string)

	// Logout
	w = s.makeRequest("POST", "/auth/logout", nil, token)
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// Thử dùng refresh token sau khi logout → phải bị lỗi
	w = s.makeRequest("POST", "/auth/refresh", map[string]string{
		"refresh_token": refreshToken,
	}, "")

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
}

// ===== W - ROLE TESTS =====

func (s *IntegrationTestSuite) Test15_User_CannotCreateProject() {
	// User thường thử tạo project → phải bị 403
	w := s.makeRequest("POST", "/auth/login", map[string]string{
		"email":    "test_user2@gmail.com",
		"password": "123456",
	}, "")

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	data := loginResponse["data"].(map[string]interface{})
	userToken := data["token"].(string)

	// User thường thử tạo project → phải bị 403
	w = s.makeRequest("POST", "/projects", map[string]string{
		"name":        "Test Project By User",
		"description": "Test",
	}, userToken)

	assert.Equal(s.T(), http.StatusForbidden, w.Code)
}

func (s *IntegrationTestSuite) Test16_User_CannotCreateTask() {
	// Login user2 có role user
	w := s.makeRequest("POST", "/auth/login", map[string]string{
		"email":    "test_user2@gmail.com",
		"password": "123456",
	}, "")

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	data := loginResponse["data"].(map[string]interface{})
	userToken := data["token"].(string)

	// Lấy project id
	var projectID int
	s.db.QueryRow(
		"SELECT id FROM projects WHERE name = 'Test Project' LIMIT 1",
	).Scan(&projectID)

	// User thường thử tạo task → phải bị 403
	w = s.makeRequest("POST", "/tasks", map[string]interface{}{
		"project_id":  projectID,
		"title":       "Test Task By User",
		"description": "Test",
		"status":      "todo",
	}, userToken)

	assert.Equal(s.T(), http.StatusForbidden, w.Code)
}

func (s *IntegrationTestSuite) Test17_Admin_CanGetAllUsers() {
	// Đổi user test thành admin
	s.db.Exec(
		"UPDATE users SET role = 'admin' WHERE email = 'test_user@gmail.com'",
	)

	// Login lại để lấy token mới có role admin
	w := s.makeRequest("POST", "/auth/login", map[string]string{
		"email":    "test_user@gmail.com",
		"password": "123456",
	}, "")

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	data := loginResponse["data"].(map[string]interface{})
	adminToken := data["token"].(string)

	// Admin lấy danh sách user
	w = s.makeRequest("GET", "/users", nil, adminToken)
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// Đổi lại thành user thường sau khi test
	s.db.Exec(
		"UPDATE users SET role = 'user' WHERE email = 'test_user@gmail.com'",
	)
}

func (s *IntegrationTestSuite) Test18_User_CannotGetAllUsers() {
	// User thường thử lấy danh sách user → phải bị 403
	w := s.makeRequest("GET", "/users", nil, s.token)
	assert.Equal(s.T(), http.StatusForbidden, w.Code)
}

func (s *IntegrationTestSuite) Test19_Admin_CanUpdateRole() {
	// Đổi user test thành admin
	s.db.Exec(
		"UPDATE users SET role = 'admin' WHERE email = 'test_user@gmail.com'",
	)

	// Login lại lấy token admin
	w := s.makeRequest("POST", "/auth/login", map[string]string{
		"email":    "test_user@gmail.com",
		"password": "123456",
	}, "")

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	data := loginResponse["data"].(map[string]interface{})
	adminToken := data["token"].(string)

	// Lấy id của user 2
	var userID int
	s.db.QueryRow(
		"SELECT id FROM users WHERE email = 'test_user2@gmail.com'",
	).Scan(&userID)

	// Admin đổi role user 2 thành manager
	w = s.makeRequest(
		"PUT",
		fmt.Sprintf("/users/%d/role", userID),
		map[string]string{
			"role": "manager",
		},
		adminToken,
	)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	// Kiểm tra role đã đổi chưa
	var role string
	s.db.QueryRow(
		"SELECT role FROM users WHERE email = 'test_user2@gmail.com'",
	).Scan(&role)

	assert.Equal(s.T(), "manager", role)

	// Đổi lại user test về user thường
	s.db.Exec(
		"UPDATE users SET role = 'user' WHERE email = 'test_user@gmail.com'",
	)
}
