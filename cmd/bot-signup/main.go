package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-go-golems/bot-signup/internal/database"
	"github.com/go-go-golems/bot-signup/internal/server"
	"github.com/go-go-golems/bot-signup/internal/web"
	"github.com/spf13/cobra"
)

var (
	addr                string
	dbPath              string
	sessionSecret       string
	discordClientID     string
	discordClientSecret string
	discordRedirectURL  string
	secureCookies       bool
	version             = "dev"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "bot-signup",
		Short: "Discord bot vibe-coding signup platform",
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(addr)
		},
	}
	serveCmd.Flags().StringVar(&addr, "addr", envOrDefault("ADDR", ":8080"), "HTTP listen address")
	serveCmd.Flags().StringVar(&dbPath, "db", envOrDefault("DB_PATH", "data/bot-signup.db"), "SQLite database path")
	serveCmd.Flags().StringVar(&sessionSecret, "session-secret", envOrDefault("SESSION_SECRET", "dev-insecure-change-me"), "session signing secret")
	serveCmd.Flags().StringVar(&discordClientID, "discord-client-id", envOrDefault("DISCORD_CLIENT_ID", ""), "Discord OAuth client ID")
	serveCmd.Flags().StringVar(&discordClientSecret, "discord-client-secret", envOrDefault("DISCORD_CLIENT_SECRET", ""), "Discord OAuth client secret")
	serveCmd.Flags().StringVar(&discordRedirectURL, "discord-redirect-url", envOrDefault("DISCORD_REDIRECT_URL", "http://localhost:8080/auth/discord/callback"), "Discord OAuth redirect URL")
	serveCmd.Flags().BoolVar(&secureCookies, "secure-cookies", envOrDefault("SECURE_COOKIES", "") == "true", "set Secure on auth cookies")

	rootCmd.AddCommand(serveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runServe(addr string) error {
	ctx := context.Background()
	db, err := database.Open(ctx, dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	mux := http.NewServeMux()
	srv := server.New(db, server.Options{
		Version:             version,
		SessionSecret:       []byte(sessionSecret),
		SecureCookies:       secureCookies,
		DiscordClientID:     discordClientID,
		DiscordClientSecret: discordClientSecret,
		DiscordRedirectURL:  discordRedirectURL,
	})
	srv.RegisterRoutes(mux)
	spaHandler, err := web.NewSPAHandler(&web.SPAOptions{APIPrefixes: []string{"/api", "/auth"}})
	if err != nil {
		log.Printf("SPA assets unavailable: %v", err)
	} else {
		mux.Handle("GET /", spaHandler)
		mux.Handle("GET /{filepath...}", spaHandler)
	}

	log.Printf("bot-signup server listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
