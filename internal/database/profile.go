package database

import (
	"context"
	"fmt"
)

func (db *DB) UpdateUserProfile(ctx context.Context, id int64, email, displayName string) (*User, error) {
	res, err := db.db.ExecContext(ctx, `
		UPDATE users
		SET email = NULLIF(?, ''), display_name = ?, updated_at = datetime('now')
		WHERE id = ?`, email, displayName, id)
	if err != nil {
		return nil, fmt.Errorf("update user profile: %w", err)
	}
	if err := requireAffected(res); err != nil {
		return nil, err
	}
	return db.GetUserByID(ctx, id)
}

func (db *DB) UpdateUserRole(ctx context.Context, id int64, role UserRole) error {
	res, err := db.db.ExecContext(ctx, `UPDATE users SET role = ?, updated_at = datetime('now') WHERE id = ?`, role, id)
	if err != nil {
		return fmt.Errorf("update user role: %w", err)
	}
	return requireAffected(res)
}
