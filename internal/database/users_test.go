package database

import (
	"context"
	"errors"
	"testing"
)

func TestUserCRUD(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)

	user, err := db.CreateUser(ctx, "123456789", "user@example.com", "CoolBotDev", "hash")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if user.ID == 0 || user.Status != UserStatusWaiting || user.Role != UserRoleUser {
		t.Fatalf("unexpected created user: %+v", user)
	}

	byEmail, err := db.GetUserByEmail(ctx, "user@example.com")
	if err != nil {
		t.Fatalf("get by email: %v", err)
	}
	if byEmail.DiscordID != "123456789" {
		t.Fatalf("unexpected discord id: %s", byEmail.DiscordID)
	}

	byDiscord, err := db.GetUserByDiscordID(ctx, "123456789")
	if err != nil {
		t.Fatalf("get by discord: %v", err)
	}
	if byDiscord.Email != "user@example.com" {
		t.Fatalf("unexpected email: %s", byDiscord.Email)
	}

	if err := db.UpdateUserStatus(ctx, user.ID, UserStatusApproved); err != nil {
		t.Fatalf("update status: %v", err)
	}
	updated, err := db.GetUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("get updated: %v", err)
	}
	if updated.Status != UserStatusApproved {
		t.Fatalf("expected approved, got %s", updated.Status)
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

func TestCreateUserUniqueness(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	if _, err := db.CreateUser(ctx, "123", "user@example.com", "First", "hash"); err != nil {
		t.Fatalf("create first user: %v", err)
	}
	if _, err := db.CreateUser(ctx, "123", "other@example.com", "Second", "hash"); err == nil {
		t.Fatal("expected duplicate discord_id error")
	}
	if _, err := db.CreateUser(ctx, "456", "user@example.com", "Third", "hash"); err == nil {
		t.Fatal("expected duplicate email error")
	}
}
