package repository

import (
	"context"
	"database/sql"
	"errors"
	"rest-api/internal/resource/model"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByCodAndEmail(cod int, email string) (*model.User, error) {
	query := `
	SELECT id, created_at, name, phone, email,cod, password_hash, activated, version
	FROM users
	WHERE email = $1 AND deleted = false AND cod = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user model.User

	err := r.db.QueryRowContext(ctx, query, email, cod).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Cod,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ID int64) (*model.User, error) {
	query := `
	SELECT id, created_at, name, phone, email, cod, password_hash, activated, version
	FROM users
	WHERE email = %1 AND deleted = false
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user model.User

	err := r.db.QueryRowContext(ctx, query, ID).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Cod,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) Insert(user *model.User) error {
	query := `
	INSERT INTO users (name, email, phone,cod, password_hash, activated,deleted)
	VALUES ($1, $2, $3, $4, $5, $6,false)
	RETURNING id, created_at, version
	`

	args := []any{
		user.Name,
		user.Email,
		user.Phone,
		user.Cod,
		user.Password.Hash,
		user.Activated,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Version,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return model.ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	query := `
	SELECT id, created_at, name, phone, email, cod, password_hash, activated, version
	FROM users
	WHERE email = $1 AND deleted = false
	`

	var user model.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Cod,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) UpdateCodByEmail(user *model.User) error {
	query := `
	UPDATE users SET
	cod = $1
	WHERE id = $2 AND version = $3
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, user.ID, user.Cod, user.Version).Scan(
		&user.Version,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return model.ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil

}

func (r *UserRepository) Update(user *model.User) error {
	query := `
	UPDATE users SET 
	name = $1, email = $2, cod = $3, phone = $4, password_hash = $5,
	activated = $6,version = version + 1
	WHERE id = $7 AND version = $8
	RETURNING version`

	args := []any{
		user.Name,
		user.Email,
		user.Cod,
		user.Phone,
		user.Password.Hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.Version,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return model.ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (r *UserRepository) Delete(user *model.User) error {
	query := `
	UPDATE users set
	deleted = true
	where id = $1 AND version = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, user.ID, user.Version)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
