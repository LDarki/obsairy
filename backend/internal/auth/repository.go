package auth

import (
	"context"
	"database/sql"

	"obsairy/internal/db"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, u *User) error

	// Sessions
	CreateSession(ctx context.Context, s *Session) error
	GetSessionByID(ctx context.Context, sessionID string) (*Session, error)
	GetSessionsByUserID(ctx context.Context, userID string) ([]Session, error)
	FindSessionByToken(ctx context.Context, token string) (any, error)
	DeleteSession(ctx context.Context, sessionID string) error
	DeleteUserSessions(ctx context.Context, userID string) error
}

// -----------------
// Postgres repo
// -----------------

type PostgresRepo struct{}

func NewRepo() *PostgresRepo {
	return &PostgresRepo{}
}

func (r *PostgresRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	row := db.Instance.QueryRowContext(ctx, `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`, email)

	var u User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresRepo) FindByID(ctx context.Context, id string) (*User, error) {
	row := db.Instance.QueryRowContext(ctx, `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`, id)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresRepo) Create(ctx context.Context, u *User) error {
	_, err := db.Instance.ExecContext(ctx, `
		INSERT INTO users (id, email, password_hash)
		VALUES ($1, $2, $3)
	`, u.ID, u.Email, u.Password)
	return err
}

func (r *PostgresRepo) CreateSession(ctx context.Context, s *Session) error {
	_, err := db.Instance.ExecContext(ctx, `
		INSERT INTO sessions (id, user_id, token, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, s.ID, s.UserID, s.Token, s.UserAgent, s.IPAddress, s.ExpiresAt)
	return err
}

func (r *PostgresRepo) GetSessionByID(ctx context.Context, sessionID string) (*Session, error) {
	row := db.Instance.QueryRowContext(ctx, `
		SELECT id, user_id, token, user_agent, ip_address, expires_at, created_at
		FROM sessions
		WHERE id = $1
	`, sessionID)

	var s Session
	err := row.Scan(&s.ID, &s.UserID, &s.Token, &s.UserAgent, &s.IPAddress, &s.ExpiresAt, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PostgresRepo) GetSessionsByUserID(ctx context.Context, userID string) ([]Session, error) {
	rows, err := db.Instance.QueryContext(ctx, `
		SELECT id, user_id, token, user_agent, ip_address, expires_at, created_at
		FROM sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions = make([]Session, 0)
	for rows.Next() {
		var s Session
		err := rows.Scan(&s.ID, &s.UserID, &s.Token, &s.UserAgent, &s.IPAddress, &s.ExpiresAt, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *PostgresRepo) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := db.Instance.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", sessionID)
	return err
}

func (r *PostgresRepo) FindSessionByToken(ctx context.Context, token string) (any, error) {
	row := db.Instance.QueryRowContext(ctx, `
		SELECT id, user_id, token, user_agent, ip_address, expires_at, created_at
		FROM sessions
		WHERE token = $1
	`, token)

	var s Session
	err := row.Scan(&s.ID, &s.UserID, &s.Token, &s.UserAgent, &s.IPAddress, &s.ExpiresAt, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PostgresRepo) DeleteUserSessions(ctx context.Context, userID string) error {
	_, err := db.Instance.ExecContext(ctx, "DELETE FROM sessions WHERE user_id = $1", userID)
	return err
}
