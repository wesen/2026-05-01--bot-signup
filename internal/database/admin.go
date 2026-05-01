package database

import (
	"context"
	"database/sql"
	"fmt"
)

type ListUsersOptions struct {
	Status  UserStatus
	Page    int
	PerPage int
}

type Stats struct {
	TotalUsers    int `json:"total_users"`
	ApprovedUsers int `json:"approved_users"`
	WaitingUsers  int `json:"waiting_users"`
	BotsRunning   int `json:"bots_running"`
}

func (db *DB) ListUsers(ctx context.Context, opts ListUsersOptions) ([]*User, int, error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PerPage < 1 {
		opts.PerPage = 20
	}
	offset := (opts.Page - 1) * opts.PerPage
	where := ""
	args := []any{}
	if opts.Status != "" {
		where = " WHERE status = ?"
		args = append(args, opts.Status)
	}

	var total int
	if err := db.db.QueryRowContext(ctx, `SELECT count(*) FROM users`+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	args = append(args, opts.PerPage, offset)
	rows, err := db.db.QueryContext(ctx, `
		SELECT id, discord_id, email, display_name, password_hash, status, role, created_at, updated_at
		FROM users`+where+`
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user, err := scanUserRows(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate users: %w", err)
	}
	return users, total, nil
}

func (db *DB) ApproveUser(ctx context.Context, userID, adminID int64, creds *BotCredentials) (*BotCredentials, error) {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin approval transaction: %w", err)
	}
	defer tx.Rollback()

	var status string
	if err := tx.QueryRowContext(ctx, `SELECT status FROM users WHERE id = ?`, userID).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user status: %w", err)
	}
	if UserStatus(status) != UserStatusWaiting {
		return nil, fmt.Errorf("user is not in waiting status")
	}
	if _, err := tx.ExecContext(ctx, `UPDATE users SET status = 'approved', updated_at = datetime('now') WHERE id = ?`, userID); err != nil {
		return nil, fmt.Errorf("approve user: %w", err)
	}
	res, err := tx.ExecContext(ctx, `
		INSERT INTO bot_credentials (user_id, application_id, bot_token, guild_id, public_key, approved_by, approved_at)
		VALUES (?, ?, ?, ?, ?, ?, datetime('now'))`, userID, creds.ApplicationID, creds.BotToken, creds.GuildID, creds.PublicKey, adminID)
	if err != nil {
		return nil, fmt.Errorf("insert approval credentials: %w", err)
	}
	credentialsID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("read credentials id: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit approval transaction: %w", err)
	}
	return db.GetCredentialsByID(ctx, credentialsID)
}

func (db *DB) GetStats(ctx context.Context) (*Stats, error) {
	stats := &Stats{}
	queries := []struct {
		target *int
		query  string
	}{
		{&stats.TotalUsers, `SELECT count(*) FROM users`},
		{&stats.ApprovedUsers, `SELECT count(*) FROM users WHERE status = 'approved'`},
		{&stats.WaitingUsers, `SELECT count(*) FROM users WHERE status = 'waiting'`},
		{&stats.BotsRunning, `SELECT count(*) FROM bot_credentials`},
	}
	for _, q := range queries {
		if err := db.db.QueryRowContext(ctx, q.query).Scan(q.target); err != nil {
			return nil, fmt.Errorf("get stats: %w", err)
		}
	}
	return stats, nil
}
