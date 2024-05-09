package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/egreb/boilerplate/internalerrors"
	"github.com/egreb/boilerplate/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionsRepository struct {
	db *pgxpool.Pool
}

func NewSessionsRespository(db *pgxpool.Pool) *SessionsRepository {
	return &SessionsRepository{
		db: db,
	}
}

func (sr *SessionsRepository) CreateSession(ctx context.Context, userId string) (*models.Session, error) {
	var s models.Session
	row := sr.db.QueryRow(ctx, `
		INSERT INTO 
			sessions (value, expires)
		VALUES ($1, $2)
		RETURNING
			id, expires`, userId, time.Now().Add(time.Hour*24))
	err := row.Scan(&s.ID, &s.Expires)
	if err != nil {
		return nil, &internalerrors.InternalError{
			Err: fmt.Errorf("unable to create session: %w", err),
		}
	}

	return &s, nil
}

func (sr *SessionsRepository) GetSession(ctx context.Context, sessionToken string) (*models.Session, error) {
	var s models.Session
	row := sr.db.QueryRow(ctx, `
		SELECT 
			id, expires
		FROM
			sessions
		WHERE
			id = $1`,
		sessionToken)
	err := row.Scan(&s.ID, &s.Expires)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, &internalerrors.NotFound{
				Err: fmt.Errorf("session not found"),
			}
		}
		return nil, &internalerrors.InternalError{
			Err: fmt.Errorf("failed getting session token: %w", err),
		}
	}

	return &s, nil
}

func (sr *SessionsRepository) Delete(ctx context.Context, sessionToken string) error {
	_, err := sr.db.Exec(ctx, `
		DELETE FROM 
			sessions
		WHERE 
			id = $1;`,
		sessionToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return &internalerrors.NotFound{
				Err: fmt.Errorf("session not found"),
			}
		}

		return &internalerrors.InternalError{
			Err: fmt.Errorf("unable to delete session: %w", err),
		}
	}

	return nil
}
