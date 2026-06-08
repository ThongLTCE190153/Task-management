package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"trithong.com/task-golang/internal/entities"
	"trithong.com/task-golang/internal/dto"
)

type TaskRepository interface {
	Create(task entities.Task, ownerID int) (entities.Task, error)
	GetAllByOwnerID(ownerID int) ([]entities.Task, error)
	GetAllByProjectID(projectID int, ownerID int) ([]entities.Task, error)
	GetAllWithFilter(ownerID int, params dto.TaskQueryParams) ([]entities.Task, int, error) // ← thêm
	GetByIDAndOwnerID(id int, ownerID int) (entities.Task, error)
	Update(id int, ownerID int, task entities.Task) (entities.Task, error)
	Delete(id int, ownerID int) error
	GetOwnerIDByProjectID(projectID int) (int, error)
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (r *taskRepository) Create(task entities.Task, ownerID int) (entities.Task, error) {

	// Bước 1: Kiểm tra project có tồn tại không
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1)`
	err := r.db.QueryRow(checkQuery, task.ProjectID).Scan(&exists)

	if err != nil {
		return entities.Task{}, err
	}

	if !exists {
		return entities.Task{}, errors.New("project not found")
	}

	// Bước 2: Kiểm tra user có phải owner không
	var isOwner bool
	ownerQuery := `SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND owner_id = $2)`
	err = r.db.QueryRow(ownerQuery, task.ProjectID, ownerID).Scan(&isOwner)

	if err != nil {
		return entities.Task{}, err
	}

	if !isOwner {
		return entities.Task{}, errors.New("forbidden")
	}

	// Bước 3: Insert task
	query := `
		INSERT INTO tasks (
			project_id,
			title,
			description,
			status,
			assignee_id
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			project_id,
			title,
			description,
			status,
			assignee_id,
			created_at
	`

	err = r.db.QueryRow(
		query,
		task.ProjectID,
		task.Title,
		task.Description,
		task.Status,
		task.AssigneeID,
	).Scan(
		&task.ID,
		&task.ProjectID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.AssigneeID,
		&task.CreatedAt,
	)

	if err != nil {
		return entities.Task{}, err
	}

	return task, nil
}

func (r *taskRepository) GetAllByOwnerID(ownerID int) ([]entities.Task, error) {
	query := `
		SELECT
			t.id,
			t.project_id,
			t.title,
			t.description,
			t.status,
			t.assignee_id,
			t.created_at
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE p.owner_id = $1
	`

	rows, err := r.db.Query(query, ownerID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []entities.Task

	for rows.Next() {
		var task entities.Task

		err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.AssigneeID,
			&task.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepository) GetByIDAndOwnerID(id int, ownerID int) (entities.Task, error) {
	var task entities.Task

	query := `
		SELECT
			t.id,
			t.project_id,
			t.title,
			t.description,
			t.status,
			t.assignee_id,
			t.created_at
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE t.id = $1
		AND p.owner_id = $2
	`

	err := r.db.QueryRow(query, id, ownerID).Scan(
		&task.ID,
		&task.ProjectID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.AssigneeID,
		&task.CreatedAt,
	)

	if err != nil {
		return entities.Task{}, err
	}

	return task, nil
}

func (r *taskRepository) Update(id int, ownerID int, task entities.Task) (entities.Task, error) {
	query := `
		UPDATE tasks t
		SET
			title = $1,
			description = $2,
			status = $3,
			assignee_id = $4
		FROM projects p
		WHERE t.project_id = p.id
		AND t.id = $5
		AND p.owner_id = $6
		RETURNING
			t.id,
			t.project_id,
			t.title,
			t.description,
			t.status,
			t.assignee_id,
			t.created_at
	`

	err := r.db.QueryRow(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.AssigneeID,
		id,
		ownerID,
	).Scan(
		&task.ID,
		&task.ProjectID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.AssigneeID,
		&task.CreatedAt,
	)

	if err != nil {
		return entities.Task{}, err
	}

	return task, nil
}

func (r *taskRepository) Delete(id int, ownerID int) error {
	query := `
		DELETE FROM tasks t
		USING projects p
		WHERE t.project_id = p.id
		AND t.id = $1
		AND p.owner_id = $2
	`

	result, err := r.db.Exec(query, id, ownerID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *taskRepository) GetOwnerIDByProjectID(projectID int) (int, error) {
	var ownerID int

	query := `
		SELECT owner_id
		FROM projects
		WHERE id = $1
	`

	err := r.db.QueryRow(query, projectID).Scan(&ownerID)

	if err != nil {
		return 0, err
	}

	return ownerID, nil
}

func (r *taskRepository) GetAllByProjectID(projectID int, ownerID int) ([]entities.Task, error) {
	query := `
		SELECT
			t.id,
			t.project_id,
			t.title,
			t.description,
			t.status,
			t.assignee_id,
			t.created_at
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE t.project_id = $1
		AND p.owner_id = $2
	`

	rows, err := r.db.Query(query, projectID, ownerID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []entities.Task

	for rows.Next() {
		var task entities.Task

		err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.AssigneeID,
			&task.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepository) GetAllWithFilter(ownerID int, params dto.TaskQueryParams) ([]entities.Task, int, error) {
	// Xử lý default value
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	// Build query động
	query := `
		SELECT
			t.id,
			t.project_id,
			t.title,
			t.description,
			t.status,
			t.assignee_id,
			t.created_at
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE p.owner_id = $1
	`

	countQuery := `
		SELECT COUNT(*)
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE p.owner_id = $1
	`

	args := []interface{}{ownerID}
	argIndex := 2

	// Thêm filter động
	if params.Status != "" {
		query += fmt.Sprintf(" AND t.status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND t.status = $%d", argIndex)
		args = append(args, params.Status)
		argIndex++
	}

	if params.AssigneeID != nil {
		query += fmt.Sprintf(" AND t.assignee_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND t.assignee_id = $%d", argIndex)
		args = append(args, *params.AssigneeID)
		argIndex++
	}

	if params.ProjectID != nil {
		query += fmt.Sprintf(" AND t.project_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND t.project_id = $%d", argIndex)
		args = append(args, *params.ProjectID)
		argIndex++
	}

	// Đếm tổng số record
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Thêm pagination
	query += fmt.Sprintf(" ORDER BY t.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []entities.Task
	for rows.Next() {
		var task entities.Task
		err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.AssigneeID,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}