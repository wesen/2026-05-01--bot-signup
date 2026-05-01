package database

// UserStatus represents the account approval state.
type UserStatus string

const (
	UserStatusWaiting   UserStatus = "waiting"
	UserStatusApproved  UserStatus = "approved"
	UserStatusRejected  UserStatus = "rejected"
	UserStatusSuspended UserStatus = "suspended"
)

// UserRole represents whether the account is a regular user or admin.
type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

// User is the database representation of a Discord-authenticated signup account.
type User struct {
	ID          int64      `json:"id"`
	DiscordID   string     `json:"discord_id"`
	Email       string     `json:"email,omitempty"`
	DisplayName string     `json:"display_name"`
	AvatarURL   string     `json:"avatar_url,omitempty"`
	Status      UserStatus `json:"status"`
	Role        UserRole   `json:"role"`
	LastLoginAt string     `json:"last_login_at,omitempty"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

// BotCredentials are the Discord application values an admin assigns on approval.
type BotCredentials struct {
	ID            int64  `json:"id"`
	UserID        int64  `json:"user_id"`
	ApplicationID string `json:"application_id"`
	BotToken      string `json:"bot_token"`
	GuildID       string `json:"guild_id"`
	PublicKey     string `json:"public_key"`
	ApprovedBy    *int64 `json:"approved_by,omitempty"`
	ApprovedAt    string `json:"approved_at,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
