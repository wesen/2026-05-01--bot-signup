package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

func (db *DB) CreateUser(ctx context.Context, discordID, email, displayName, passwordHash string) (*User, error) {
	res, err := db.db.ExecContext(ctx, `
		INSERT INTO users (discord_id, email, display_name, password_hash)
		VALUES (?, ?, ?, ?)`, discordID, email, displayName, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("read created user id: %w", err)
	}
	return db.GetUserByID(ctx, id)
}

func (db *DB) GetUserByID(ctx context.Context, id int64) (*User, error) {
	return db.scanUser(db.db.QueryRowContext(ctx, `
		SELECT id, discord_id, email, display_name, password_hash, status, role, created_at, updated_at
		FROM users WHERE id = ?`, id))
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return db.scanUser(db.db.QueryRowContext(ctx, `
		SELECT id, discord_id, email, display_name, password_hash, status, role, created_at, updated_at
		FROM users WHERE email = ?`, email))
}

func (db *DB) GetUserByDiscordID(ctx context.Context, discordID string) (*User, error) {
	return db.scanUser(db.db.QueryRowContext(ctx, `
		SELECT id, discord_id, email, display_name, password_hash, status, role, created_at, updated_at
		FROM users WHERE discord_id = ?`, discordID))
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

	rows, err := db.db.QueryContext(ctx, `
		SELECT id, discord_id, email, display_name, password_hash, status, role, created_at, updated_at
		FROM users
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

func (db *DB) scanUser(row *sql.Row) (*User, error) {
	user := &User{}
	var status string
	var role string
	err := row.Scan(&user.ID, &user.DiscordID, &user.Email, &user.DisplayName, &user.PasswordHash, &status, &role, &user.CreatedAt, &user.UpdatedAt)
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
	if err := row.Scan(&user.ID, &user.DiscordID, &user.Email, &user.DisplayName, &user.PasswordHash, &status, &role, &user.CreatedAt, &user.UpdatedAt); err != nil {
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
