package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-go-golems/bot-signup/internal/database"
	"github.com/go-go-golems/bot-signup/internal/server"
	"github.com/spf13/cobra"
)

var (
	addr    string
	dbPath  string
	version = "dev"
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
	srv := server.New(db, version)
	srv.RegisterRoutes(mux)

	log.Printf("bot-signup server listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
