package database

import (
	"context"
	"errors"
	"testing"
)

func TestUserCRUD(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)

	user, err := db.UpsertDiscordUser(ctx, "123456789", "user@example.com", "CoolBotDev", "https://cdn.example/avatar.png")
	if err != nil {
		t.Fatalf("upsert user: %v", err)
	}
	if user.ID == 0 || user.Status != UserStatusWaiting || user.Role != UserRoleUser {
		t.Fatalf("unexpected created user: %+v", user)
	}
	if user.AvatarURL == "" || user.LastLoginAt == "" {
		t.Fatalf("expected avatar and last login to be set: %+v", user)
	}

	byEmail, err := db.GetUserByEmail(ctx, "user@example.com")
	if err != nil {
		t.Fatalf("get by email: %v", err)
	}
	if byEmail.DiscordID != "123456789" {
		t.Fatalf("unexpected discord id: %s", byEmail.DiscordID)
	}

	updatedOAuth, err := db.UpsertDiscordUser(ctx, "123456789", "new@example.com", "NewName", "")
	if err != nil {
		t.Fatalf("update existing oauth user: %v", err)
	}
	if updatedOAuth.ID != user.ID || updatedOAuth.Email != "new@example.com" || updatedOAuth.DisplayName != "NewName" {
		t.Fatalf("unexpected oauth update: %+v", updatedOAuth)
	}

	if err := db.UpdateUserStatus(ctx, user.ID, UserStatusApproved); err != nil {
		t.Fatalf("update status: %v", err)
	}
	users, total, err := db.ListUsersByStatus(ctx, UserStatusApproved, 1, 20)
	if err != nil {
		t.Fatalf("list approved: %v", err)
	}
	if total != 1 || len(users) != 1 {
		t.Fatalf("expected one approved user, total=%d len=%d", total, len(users))
	}

	if err := db.DeleteUser(ctx, user.ID); err != nil {
		t.Fatalf("delete user: %v", err)
	}
	if _, err := db.GetUserByID(ctx, user.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestUpsertDiscordUserUniqueEmail(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	if _, err := db.UpsertDiscordUser(ctx, "123", "user@example.com", "First", ""); err != nil {
		t.Fatalf("create first user: %v", err)
	}
	if _, err := db.UpsertDiscordUser(ctx, "456", "user@example.com", "Second", ""); err == nil {
		t.Fatal("expected duplicate email error")
	}
}
