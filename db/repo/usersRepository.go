package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/egreb/boilerplate/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

func NewUsersRepository(db *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{db}
}

func (u *UsersRepository) CreateUser(ctx context.Context, username string, password string, salt []byte) error {
	_, err := u.db.Exec(ctx, `
		INSERT INTO 
			users (username, password, salt) 
		VALUES ($1, $2, $3)
		ON CONFLICT
			DO NOTHING;
	`, username, password, salt)
	if err != nil {
		if perr, ok := err.(*pq.Error); ok {
			if perr.Constraint != "" {
				return errors.BadRequest{Err: fmt.Errorf("unable to create user because of constraint: %w", err)}
			}
		}

		return errors.InternalError{Err: fmt.Errorf("failed to create user: %w", err)}
	}

	return nil
}

func (u *UsersRepository) GetUserByUsername(ctx context.Context, username string) (string, string, []byte, error) {
	var id string
	var password string
	var salt []byte

	row := u.db.QueryRow(ctx, `
		SELECT 
			id,
			password,
			salt
		FROM
			users
		WHERE
			username = $1`,
		username)

	err := row.Scan(&id, &password, &salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return id, password, salt, errors.NotFound{
				Err: fmt.Errorf("username %s not found: %w", username, err),
			}
		}

		return id, password, salt, errors.InternalError{
			Err: fmt.Errorf("failed getting user by username: %w", err),
		}
	}

	return id, password, salt, nil
}