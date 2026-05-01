package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

func (db *DB) UpsertDiscordUser(ctx context.Context, discordID, email, displayName, avatarURL string) (*User, error) {
	if displayName == "" {
		displayName = "Discord User"
	}
	_, err := db.db.ExecContext(ctx, `
		INSERT INTO users (discord_id, email, display_name, avatar_url, last_login_at)
		VALUES (?, NULLIF(?, ''), ?, NULLIF(?, ''), datetime('now'))
		ON CONFLICT(discord_id) DO UPDATE SET
			email = COALESCE(NULLIF(excluded.email, ''), users.email),
			display_name = excluded.display_name,
			avatar_url = COALESCE(NULLIF(excluded.avatar_url, ''), users.avatar_url),
			last_login_at = datetime('now'),
			updated_at = datetime('now')`, discordID, email, displayName, avatarURL)
	if err != nil {
		return nil, fmt.Errorf("upsert discord user: %w", err)
	}
	return db.GetUserByDiscordID(ctx, discordID)
}

func (db *DB) GetUserByID(ctx context.Context, id int64) (*User, error) {
	return db.scanUser(db.db.QueryRowContext(ctx, selectUserSQL+` WHERE id = ?`, id))
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return db.scanUser(db.db.QueryRowContext(ctx, selectUserSQL+` WHERE email = ?`, email))
}

func (db *DB) GetUserByDiscordID(ctx context.Context, discordID string) (*User, error) {
	return db.scanUser(db.db.QueryRowContext(ctx, selectUserSQL+` WHERE discord_id = ?`, discordID))
}

func (db *DB) UpdateUserStatus(ctx context.Context, id int64, status UserStatus) error {
	res, err := db.db.ExecContext(ctx, `UPDATE users SET status = ?, updated_at = datetime('now') WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("update user status: %w", err)
	}
	return requireAffected(res)
}

func (db *DB) ListUsersByStatus(ctx context.Context, status UserStatus, page, perPage int) ([]*User, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	var total int
	if err := db.db.QueryRowContext(ctx, `SELECT count(*) FROM users WHERE status = ?`, status).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users by status: %w", err)
	}

	rows, err := db.db.QueryContext(ctx, selectUserSQL+`
		WHERE status = ?
		ORDER BY created_at ASC, id ASC
		LIMIT ? OFFSET ?`, status, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list users by status: %w", err)
	}
	defer rows.Close()

	users := make([]*User, 0)
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

func (db *DB) DeleteUser(ctx context.Context, id int64) error {
	res, err := db.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return requireAffected(res)
}

const selectUserSQL = `
	SELECT id, discord_id, COALESCE(email, ''), display_name, COALESCE(avatar_url, ''), status, role, COALESCE(last_login_at, ''), created_at, updated_at
	FROM users`

func (db *DB) scanUser(row *sql.Row) (*User, error) {
	user := &User{}
	var status string
	var role string
	err := row.Scan(&user.ID, &user.DiscordID, &user.Email, &user.DisplayName, &user.AvatarURL, &status, &role, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	user.Status = UserStatus(status)
	user.Role = UserRole(role)
	return user, nil
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUserRows(row userScanner) (*User, error) {
	user := &User{}
	var status string
	var role string
	if err := row.Scan(&user.ID, &user.DiscordID, &user.Email, &user.DisplayName, &user.AvatarURL, &status, &role, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, fmt.Errorf("scan user row: %w", err)
	}
	user.Status = UserStatus(status)
	user.Role = UserRole(role)
	return user, nil
}

func requireAffected(res sql.Result) error {
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
