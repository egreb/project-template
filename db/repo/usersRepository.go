package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/egreb/boilerplate/internalerrors"
	"github.com/egreb/boilerplate/models"
	"github.com/jackc/pgx/v5"
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
				return &internalerrors.BadRequest{Err: fmt.Errorf("unable to create user because of constraint: %w", err)}
			}
		}

		return &internalerrors.InternalError{Err: fmt.Errorf("failed to create user: %w", err)}
	}

	return nil
}

func (u *UsersRepository) GetUserCredentialsByUsername(ctx context.Context, username string) (string, string, []byte, error) {
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
		if err == pgx.ErrNoRows {
			return id, password, salt, &internalerrors.NotFound{
				Err: fmt.Errorf("user by %s not found", username),
			}
		}
		return id, password, salt, &internalerrors.InternalError{
			Err: fmt.Errorf("failed getting user by username: %w", err),
		}
	}

	return id, password, salt, nil
}

func (u *UsersRepository) GetMe(ctx context.Context, sessionToken string) (*models.Me, error) {
	var me models.Me

	row := u.db.QueryRow(ctx, `
	SELECT 
		id, 
		username
	FROM 
		users
	WHERE
		id = (
			SELECT 
				value
			FROM
				sessions
			WHERE
				id = $1
		)::uuid;
	`, sessionToken)
	err := row.Scan(&me.ID, &me.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &internalerrors.NotFound{}
		}

		return nil, &internalerrors.InternalError{
			Err: fmt.Errorf("something went wrong looking up user: %w", err),
		}
	}

	return &me, nil
}
