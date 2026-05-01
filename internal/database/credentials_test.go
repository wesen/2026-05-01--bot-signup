package database

import (
	"context"
	"errors"
	"testing"
)

func TestBotCredentialsCRUD(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	user, err := db.CreateUser(ctx, "123", "user@example.com", "User", "hash")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	admin, err := db.CreateUser(ctx, "999", "admin@example.com", "Admin", "hash")
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}

	created, err := db.InsertBotCredentials(ctx, &BotCredentials{
		UserID:        user.ID,
		ApplicationID: "987654321",
		BotToken:      "token",
		GuildID:       "111222333",
		PublicKey:     "abcdef",
		ApprovedBy:    &admin.ID,
	})
	if err != nil {
		t.Fatalf("insert credentials: %v", err)
	}
	if created.ID == 0 || created.ApprovedAt == "" || created.ApprovedBy == nil || *created.ApprovedBy != admin.ID {
		t.Fatalf("unexpected created credentials: %+v", created)
	}

	byUser, err := db.GetCredentialsByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("get credentials by user: %v", err)
	}
	if byUser.ApplicationID != "987654321" {
		t.Fatalf("unexpected application id: %s", byUser.ApplicationID)
	}

	byUser.ApplicationID = "123456789"
	byUser.BotToken = "new-token"
	if err := db.UpdateBotCredentials(ctx, byUser); err != nil {
		t.Fatalf("update credentials: %v", err)
	}
	updated, err := db.GetCredentialsByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("get updated credentials: %v", err)
	}
	if updated.ApplicationID != "123456789" || updated.BotToken != "new-token" {
		t.Fatalf("unexpected updated credentials: %+v", updated)
	}

	if err := db.DeleteUser(ctx, user.ID); err != nil {
		t.Fatalf("delete user: %v", err)
	}
	if _, err := db.GetCredentialsByUserID(ctx, user.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected cascade delete, got %v", err)
	}
}
