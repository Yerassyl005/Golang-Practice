package users

import (
	"context"
	"errors"
	"time"

	"github.com/Yerassyl005/go-practice3/internal/repository/_postgres"
	"github.com/Yerassyl005/go-practice3/pkg/modules"
)

type Repository struct {
	db              *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: 5 * time.Second,
	}
}

//////////////////////
// 1️⃣ GET ALL USERS
//////////////////////

func (r *Repository) GetUsers(limit, offset int) ([]modules.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	var users []modules.User

	query := `
		SELECT id, name, email, age, city, deleted_at
		FROM users
		WHERE deleted_at IS NULL
		LIMIT $1 OFFSET $2
	`

	err := r.db.DB.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}

//////////////////////
// 2️⃣ GET USER BY ID
//////////////////////

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	var user modules.User

	query := `SELECT id, name, email, age, city FROM users WHERE id = $1`

	err := r.db.DB.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

//////////////////////
// 3️⃣ CREATE USER
//////////////////////

func (r *Repository) CreateUser(user *modules.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	query := `
		INSERT INTO users (name, email, age, city)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int

	err := r.db.DB.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Age,
		user.City,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

//////////////////////
// 4️⃣ UPDATE USER
//////////////////////

func (r *Repository) UpdateUser(user *modules.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	query := `
		UPDATE users
		SET name=$1, email=$2, age=$3, city=$4
		WHERE id=$5
	`

	result, err := r.db.DB.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Age,
		user.City,
		user.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

//////////////////////
// 5️⃣ DELETE USER
//////////////////////

func (r *Repository) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id=$1 AND deleted_at IS NULL
	`

	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}