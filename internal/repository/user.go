package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yosakoo/task-traker/internal/domain"
	"github.com/yosakoo/task-traker/internal/domain/models"
	"github.com/yosakoo/task-traker/pkg/postgres"
)

type UserRepo struct {
	s *postgres.Storage
}

func NewUserRepo(pg *postgres.Storage) *UserRepo {
	return &UserRepo{s: pg}
}

func (r *UserRepo) GetUserByCredentials(ctx context.Context, email string, password []byte) (*models.User, error) {
	var user models.User
	err := r.s.Pool.QueryRow(ctx, "SELECT id, email, pass_hash FROM users WHERE email = $1 AND pass_hash = $2", email, password).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
        fmt.Println(err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserByRefresh(ctx context.Context, refresh string) (int, error) {
	var userId int
	var expiresAt time.Time
	err := r.s.Pool.QueryRow(ctx, "SELECT user_id, expires_at FROM refresh_tokens WHERE refresh_token = $1", refresh).Scan(&userId, &expiresAt)
	if err != nil {
        
		return userId, err
	}
	if time.Now().After(expiresAt) {
		return userId, domain.ErrTokenExpired
	}

	return userId, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, email FROM users WHERE id = $1"
	err := r.s.Pool.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) AddUser(ctx context.Context, user models.User) (int, error) {
	txOptions := pgx.TxOptions{}

	tx, err := r.s.Pool.BeginTx(ctx, txOptions)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var userId int
	var existingUserId int
	err = tx.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", user.Email).Scan(&existingUserId)
	if err == nil {
		return 0, domain.ErrUserAlreadyExists
	} else if err != pgx.ErrNoRows {
		return 0, err
	}

	err = tx.QueryRow(ctx, "INSERT INTO users (email,name, pass_hash) VALUES ($1, $2,$3) RETURNING id", user.Email, user.Name, user.Password).Scan(&userId)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(ctx, "INSERT INTO refresh_tokens (user_id,refresh_token) VALUES ($1, NULL)", userId)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, errors.New("error committing database transaction")
	}

	return userId, nil
}

func (r *UserRepo) SetSession(ctx context.Context, userId int, refresh string, expiresAt time.Time) error {
	txOptions := pgx.TxOptions{}

	tx, err := r.s.Pool.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE refresh_tokens SET refresh_token = $2, expires_at = $3 WHERE user_id = $1", userId, refresh, expiresAt)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.New("error committing database transaction")
	}

	return nil
}
