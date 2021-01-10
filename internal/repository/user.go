package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zapkub/cftl/internal"
	"github.com/zapkub/cftl/internal/apperror"
	"github.com/zapkub/cftl/internal/database"
)

func getSession(ctx context.Context, db *database.DB, accessToken string) (_ *internal.Session, err error) {
	var session internal.Session
	err = db.QueryRow(
		ctx,
		`SELECT s.access_token,
				s.email,
				s.origin,
				s.refresh_token
		 FROM sessions s
		 WHERE s.access_token = $1
		`,
		accessToken,
	).Scan(&session.Email, &session.Email, &session.Origin, &session.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (repo *DB) InsertSession(ctx context.Context, s *internal.Session) (string, error) {
	return insertSession(ctx, repo.db, s)
}
func insertSession(ctx context.Context, db *database.DB, s *internal.Session) (string, error) {
	var token string
	err := db.QueryRow(
		ctx,
		`INSERT INTO sessions(
			access_token,
			email,
			origin,
			refresh_token
		)
		VALUES ($1,$2,$3,$4)
		RETURNING access_token
		`,
		s.AccessToken,
		s.Email,
		s.Origin,
		s.RefreshToken,
	).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}
func (repo *DB) GetUser(ctx context.Context, email string) (*internal.User, error) {
	userinfo, err := getUser(ctx, repo.db, email)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", apperror.NotFound, err)
		}
		return nil, err
	}
	return userinfo, nil
}

func getUser(ctx context.Context, db *database.DB, email string) (_ *internal.User, err error) {
	var user internal.User
	err = db.QueryRow(
		ctx,
		`SELECT u.email,
				u.username,
				u.name
		FROM users u
		WHERE u.email = $1
		`,
		email,
	).Scan(&user.Email, &user.Username, &user.Name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *DB) InsertUser(ctx context.Context, user *internal.User) (_ string, err error) {
	return insertUser(ctx, repo.db, user)
}

func insertUser(ctx context.Context, db *database.DB, u *internal.User) (_ string, err error) {

	var userID string
	err = db.QueryRow(ctx,
		`INSERT INTO users(
				email,
				username,
				name)
		VALUES($1,$2,$3)
		RETURNING email`,
		u.Email,
		u.Username,
		u.Name,
	).Scan(&userID)

	if err != nil {
		return "", err
	}

	return userID, nil
}
