package repositories

import (
	"database/sql"

	"trithong.com/task-golang/internal/entities"
	"trithong.com/task-golang/internal/dto"
)

type ProjectRepository interface {
	Create(project entities.Project) (entities.Project, error)
	GetAllByOwnerID(ownerID int) ([]entities.Project, error)
	GetAllWithPagination(ownerID int, params dto.ProjectQueryParams) ([]entities.Project, int, error)
	GetByIDAndOwnerID(id int, ownerID int) (entities.Project, error)
	Update(id int, ownerID int, project entities.Project) (entities.Project, error)
	Delete(id int, ownerID int) error
}

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{
		db: db,
	}
}

func (r *projectRepository) Create(project entities.Project) (entities.Project, error) {
	query := `
		INSERT INTO projects (name, description, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, owner_id, created_at
	`

	err := r.db.QueryRow(
		query,
		project.Name,
		project.Description,
		project.OwnerID,
	).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.OwnerID,
		&project.CreatedAt,
	)

	if err != nil {
		return entities.Project{}, err
	}

	return project, nil
}

func (r *projectRepository) GetAllByOwnerID(ownerID int) ([]entities.Project, error) {
	query := `
		SELECT id, name, description, owner_id, created_at
		FROM projects
		WHERE owner_id = $1
	`

	rows, err := r.db.Query(query, ownerID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var projects []entities.Project

	for rows.Next() {
		var project entities.Project

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.OwnerID,
			&project.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (r *projectRepository) GetByIDAndOwnerID(id int, ownerID int) (entities.Project, error) {
	var project entities.Project

	query := `
		SELECT id, name, description, owner_id, created_at
		FROM projects
		WHERE id = $1 AND owner_id = $2
	`

	err := r.db.QueryRow(query, id, ownerID).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.OwnerID,
		&project.CreatedAt,
	)

	if err != nil {
		return entities.Project{}, err
	}

	return project, nil
}

func (r *projectRepository) Update(id int, ownerID int, project entities.Project) (entities.Project, error) {
	query := `
		UPDATE projects
		SET name = $1, description = $2
		WHERE id = $3 AND owner_id = $4
		RETURNING id, name, description, owner_id, created_at
	`

	err := r.db.QueryRow(
		query,
		project.Name,
		project.Description,
		id,
		ownerID,
	).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.OwnerID,
		&project.CreatedAt,
	)

	if err != nil {
		return entities.Project{}, err
	}

	return project, nil
}

func (r *projectRepository) Delete(id int, ownerID int) error {
	query := `
		DELETE FROM projects
		WHERE id = $1 AND owner_id = $2
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

func (r *projectRepository) GetAllWithPagination(ownerID int, params dto.ProjectQueryParams) ([]entities.Project, int, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	// Đếm tổng
	var total int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM projects WHERE owner_id = $1",
		ownerID,
	).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	// Lấy data
	query := `
		SELECT id, name, description, owner_id, created_at
		FROM projects
		WHERE owner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, ownerID, params.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []entities.Project
	for rows.Next() {
		var project entities.Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.OwnerID,
			&project.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		projects = append(projects, project)
	}

	return projects, total, nil
}