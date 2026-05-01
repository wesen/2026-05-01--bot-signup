package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

var discordEndpoint = oauth2.Endpoint{
	AuthURL:  "https://discord.com/oauth2/authorize",
	TokenURL: "https://discord.com/api/oauth2/token",
}

type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	GlobalName    string `json:"global_name"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
	Verified      bool   `json:"verified"`
}

func (u DiscordUser) DisplayName() string {
	if u.GlobalName != "" {
		return u.GlobalName
	}
	return u.Username
}

func (u DiscordUser) AvatarURL() string {
	if u.ID == "" || u.Avatar == "" {
		return ""
	}
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", u.ID, u.Avatar)
}

type DiscordOAuth struct {
	config *oauth2.Config
	client *http.Client
}

func NewDiscordOAuth(clientID, clientSecret, redirectURL string) *DiscordOAuth {
	return &DiscordOAuth{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"identify", "email"},
			Endpoint:     discordEndpoint,
		},
		client: http.DefaultClient,
	}
}

func (d *DiscordOAuth) AuthCodeURL(state string) string {
	return d.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (d *DiscordOAuth) ExchangeAndFetchUser(ctx context.Context, code string) (*DiscordUser, error) {
	token, err := d.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange discord code: %w", err)
	}
	client := d.config.Client(ctx, token)
	if d.client != nil && d.client.Transport != nil {
		// Preserve custom transports in tests without replacing the OAuth2 transport
		// that injects the Authorization header for Discord API calls.
		if transport, ok := client.Transport.(*oauth2.Transport); ok {
			transport.Base = d.client.Transport
		} else {
			client.Transport = d.client.Transport
		}
	}
	return fetchDiscordUser(ctx, client)
}

func fetchDiscordUser(ctx context.Context, client *http.Client) (*DiscordUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch discord user: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("discord user fetch failed: status %d", resp.StatusCode)
	}
	var user DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decode discord user: %w", err)
	}
	if user.ID == "" {
		return nil, fmt.Errorf("discord user missing id")
	}
	return &user, nil
}
