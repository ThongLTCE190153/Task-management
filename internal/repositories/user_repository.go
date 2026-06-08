package repositories

import (
	"database/sql"

	"trithong.com/task-golang/internal/entities"
)

type UserRepository interface {
	Create(user entities.User) (entities.User, error)
	GetByID(id int) (entities.User, error)
	GetByEmail(email string) (entities.User, error)
	Update(id int, user entities.User) (entities.User, error)
	Delete(id int) error
	GetAll() ([]entities.User, error)     // ← thêm
	UpdateRole(id int, role string) error // ← thêm
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user entities.User) (entities.User, error) {
	query := `
		INSERT INTO users (email, password_hash, full_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, full_name, role, created_at
	`

	err := r.db.QueryRow(
		query,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Role,
	).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetByID(id int) (entities.User, error) {
	var user entities.User

	query := `
		SELECT id, email, full_name, role, created_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetByEmail(email string) (entities.User, error) {
	var user entities.User

	query := `
		SELECT id, email, password_hash, full_name, role, created_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r *userRepository) Update(id int, user entities.User) (entities.User, error) {
	query := `
		UPDATE users
		SET full_name = $1
		WHERE id = $2
		RETURNING id, email, full_name, role, created_at
	`

	err := r.db.QueryRow(
		query,
		user.FullName,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r *userRepository) Delete(id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)

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

func (r *userRepository) GetAll() ([]entities.User, error) {
	query := `
		SELECT id, email, full_name, role, created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FullName,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) UpdateRole(id int, role string) error {
	query := `UPDATE users SET role = $1 WHERE id = $2`

	result, err := r.db.Exec(query, role, id)
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