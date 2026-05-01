package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (db *DB) InsertBotCredentials(ctx context.Context, creds *BotCredentials) (*BotCredentials, error) {
	res, err := db.db.ExecContext(ctx, `
		INSERT INTO bot_credentials (user_id, application_id, bot_token, guild_id, public_key, approved_by, approved_at)
		VALUES (?, ?, ?, ?, ?, ?, COALESCE(?, datetime('now')))`,
		creds.UserID, creds.ApplicationID, creds.BotToken, creds.GuildID, creds.PublicKey, creds.ApprovedBy, emptyToNil(creds.ApprovedAt))
	if err != nil {
		return nil, fmt.Errorf("insert bot credentials: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("read created credentials id: %w", err)
	}
	return db.GetCredentialsByID(ctx, id)
}

func (db *DB) GetCredentialsByID(ctx context.Context, id int64) (*BotCredentials, error) {
	return scanCredentials(db.db.QueryRowContext(ctx, `
		SELECT id, user_id, application_id, bot_token, guild_id, public_key, approved_by, approved_at, created_at, updated_at
		FROM bot_credentials WHERE id = ?`, id))
}

func (db *DB) GetCredentialsByUserID(ctx context.Context, userID int64) (*BotCredentials, error) {
	return scanCredentials(db.db.QueryRowContext(ctx, `
		SELECT id, user_id, application_id, bot_token, guild_id, public_key, approved_by, approved_at, created_at, updated_at
		FROM bot_credentials WHERE user_id = ?`, userID))
}

func (db *DB) UpdateBotCredentials(ctx context.Context, creds *BotCredentials) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin credential update transaction: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
		UPDATE bot_credentials
		SET application_id = ?, bot_token = ?, guild_id = ?, public_key = ?, updated_at = datetime('now')
		WHERE user_id = ?`, creds.ApplicationID, creds.BotToken, creds.GuildID, creds.PublicKey, creds.UserID)
	if err != nil {
		return fmt.Errorf("update bot credentials: %w", err)
	}
	if err := requireAffected(res); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE users SET status = 'approved', updated_at = datetime('now') WHERE id = ?`, creds.UserID); err != nil {
		return fmt.Errorf("approve user after credential update: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit credential update transaction: %w", err)
	}
	return nil
}

func scanCredentials(row *sql.Row) (*BotCredentials, error) {
	creds := &BotCredentials{}
	var approvedBy sql.NullInt64
	var approvedAt sql.NullString
	err := row.Scan(&creds.ID, &creds.UserID, &creds.ApplicationID, &creds.BotToken, &creds.GuildID, &creds.PublicKey, &approvedBy, &approvedAt, &creds.CreatedAt, &creds.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan credentials: %w", err)
	}
	if approvedBy.Valid {
		creds.ApprovedBy = &approvedBy.Int64
	}
	if approvedAt.Valid {
		creds.ApprovedAt = approvedAt.String
	}
	return creds, nil
}

func emptyToNil(s string) any {
	if s == "" {
		return nil
	}
	return s
}
