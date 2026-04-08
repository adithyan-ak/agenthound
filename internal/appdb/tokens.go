package appdb

import (
	"context"
	"fmt"

	"github.com/adithyan-ak/agenthound/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenStore struct {
	pool *pgxpool.Pool
}

func NewTokenStore(pool *pgxpool.Pool) *TokenStore {
	return &TokenStore{pool: pool}
}

func (s *TokenStore) Create(ctx context.Context, token *model.APIToken) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO api_tokens (id, user_id, token_hash, name, created_at, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		token.ID, token.UserID, token.TokenHash, token.Name, token.CreatedAt, token.ExpiresAt)
	if err != nil {
		return fmt.Errorf("create token: %w", err)
	}
	return nil
}

func (s *TokenStore) GetByHash(ctx context.Context, tokenHash string) (*model.APIToken, error) {
	token := &model.APIToken{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, token_hash, name, created_at, expires_at, last_used
		 FROM api_tokens WHERE token_hash = $1`, tokenHash).
		Scan(&token.ID, &token.UserID, &token.TokenHash, &token.Name, &token.CreatedAt,
			&token.ExpiresAt, &token.LastUsed)
	if err != nil {
		return nil, fmt.Errorf("get token by hash: %w", err)
	}
	return token, nil
}

func (s *TokenStore) UpdateLastUsed(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE api_tokens SET last_used = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("update token last_used: %w", err)
	}
	return nil
}

func (s *TokenStore) ListByUser(ctx context.Context, userID string) ([]model.APIToken, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, token_hash, name, created_at, expires_at, last_used
		 FROM api_tokens WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list tokens: %w", err)
	}
	defer rows.Close()

	var tokens []model.APIToken
	for rows.Next() {
		var t model.APIToken
		if err := rows.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.Name, &t.CreatedAt,
			&t.ExpiresAt, &t.LastUsed); err != nil {
			return nil, fmt.Errorf("scan token row: %w", err)
		}
		tokens = append(tokens, t)
	}
	return tokens, rows.Err()
}

func (s *TokenStore) Delete(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx,
		`DELETE FROM api_tokens WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete token: %w", err)
	}
	return nil
}
