package appdb

import (
	"context"
	"fmt"

	"github.com/adithyan-ak/agenthound/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	pool *pgxpool.Pool
}

func NewUserStore(pool *pgxpool.Pool) *UserStore {
	return &UserStore{pool: pool}
}

func (s *UserStore) Create(ctx context.Context, user *model.User) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO users (id, username, password_hash, role, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		user.ID, user.Username, user.PasswordHash, user.Role, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (s *UserStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, role, created_at, last_login
		 FROM users WHERE username = $1`, username).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.LastLogin)
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return user, nil
}

func (s *UserStore) GetByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, role, created_at, last_login
		 FROM users WHERE id = $1`, id).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.LastLogin)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (s *UserStore) UpdateLastLogin(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET last_login = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("update last login: %w", err)
	}
	return nil
}

func (s *UserStore) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.pool.Query(ctx,
		`SELECT id, username, password_hash, role, created_at, last_login
		 FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.LastLogin); err != nil {
			return nil, fmt.Errorf("scan user row: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (s *UserStore) Delete(ctx context.Context, id string) (bool, error) {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return false, fmt.Errorf("delete user: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

func (s *UserStore) Count(ctx context.Context) (int, error) {
	var count int
	err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return count, nil
}
