package repositories

import (
	"database/sql"

	"trithong.com/task-golang/internal/entities"
)

type CommentRepository interface {
	Create(comment entities.Comment) (entities.Comment, error)
	GetByTaskID(taskID int) ([]entities.Comment, error)
	CheckTaskBelongsToUser(taskID int, userID int) (bool, error)
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (r *commentRepository) Create(
	comment entities.Comment,
) (entities.Comment, error) {

	query := `
		INSERT INTO task_comments
		(task_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		comment.TaskID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)

	if err != nil {
		return entities.Comment{}, err
	}

	return comment, nil
}

func (r *commentRepository) GetByTaskID(
	taskID int,
) ([]entities.Comment, error) {

	query := `
		SELECT
			id,
			task_id,
			user_id,
			content,
			created_at
		FROM task_comments
		WHERE task_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, taskID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []entities.Comment

	for rows.Next() {
		var comment entities.Comment

		err := rows.Scan(
			&comment.ID,
			&comment.TaskID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *commentRepository) CheckTaskBelongsToUser(taskID int, userID int) (bool, error) {
    query := `
        SELECT COUNT(*)
        FROM tasks t
        JOIN projects p ON t.project_id = p.id
        WHERE t.id = $1
        AND p.owner_id = $2
    `

    var count int
    err := r.db.QueryRow(query, taskID, userID).Scan(&count)

    if err != nil {
        return false, err
    }

    return count > 0, nil
}