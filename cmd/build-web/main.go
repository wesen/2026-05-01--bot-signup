// build-web builds the React frontend and copies the output to
// internal/web/embed/public/ for embedding via //go:embed.
//
// Usage:
//
//	go run ./cmd/build-web                    # Dagger (requires Docker/Dagger engine)
//	BUILD_WEB_LOCAL=1 go run ./cmd/build-web  # Local pnpm (requires node + pnpm)
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

const (
	defaultBuilderImage = "node:22-bookworm"
	defaultPNPMVersion  = "10.15.1"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "build-web: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}
	if os.Getenv("BUILD_WEB_LOCAL") == "1" {
		return runLocal(repoRoot)
	}
	if err := runDagger(ctx, repoRoot); err != nil {
		if errors.Is(err, errDaggerUnavailable) {
			fmt.Fprintln(os.Stderr, "dagger unavailable, falling back to local pnpm")
			return runLocal(repoRoot)
		}
		return err
	}
	return nil
}

var errDaggerUnavailable = errors.New("dagger: engine not reachable")

func runDagger(ctx context.Context, repoRoot string) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return fmt.Errorf("%w: %v", errDaggerUnavailable, err)
	}
	defer func() { _ = client.Close() }()

	uiDir := filepath.Join(repoRoot, "ui")
	pnpmVersion := getenv("WEB_PNPM_VERSION", readPNPMVersion(filepath.Join(uiDir, "package.json")))
	if pnpmVersion == "" {
		pnpmVersion = defaultPNPMVersion
	}
	builderImage := getenv("WEB_BUILDER_IMAGE", defaultBuilderImage)

	source := client.Host().Directory(uiDir, dagger.HostDirectoryOpts{
		Exclude: []string{"dist", "node_modules", "storybook-static"},
	})

	pnpmStore := client.CacheVolume("bot-signup-ui-pnpm-store")
	pathEnv := "/pnpm:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

	container := client.Container().
		From(builderImage).
		WithEnvVariable("PNPM_HOME", "/pnpm").
		WithEnvVariable("PATH", pathEnv).
		WithMountedCache("/pnpm/store", pnpmStore).
		WithDirectory("/ui", source).
		WithWorkdir("/ui").
		WithExec([]string{"sh", "-lc", "corepack enable && corepack prepare pnpm@" + pnpmVersion + " --activate"}).
		WithExec([]string{"pnpm", "install", "--frozen-lockfile", "--prefer-offline"}).
		WithExec([]string{"pnpm", "run", "build"})

	tmpDir, err := os.MkdirTemp("", "bot-signup-ui-dist-")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	if _, err := container.Directory("/ui/dist").Export(ctx, tmpDir); err != nil {
		return fmt.Errorf("export dist: %w", err)
	}
	return copyDistToEmbed(repoRoot, tmpDir, "Dagger")
}

func runLocal(repoRoot string) error {
	if err := runCmd(repoRoot, "pnpm", "--dir", "ui", "run", "build"); err != nil {
		return fmt.Errorf("pnpm build (local): %w", err)
	}
	return copyDistToEmbed(repoRoot, filepath.Join(repoRoot, "ui", "dist"), "local pnpm")
}

func copyDistToEmbed(repoRoot, src, mode string) error {
	dst := filepath.Join(repoRoot, "internal", "web", "embed", "public")
	if err := recreate(dst); err != nil {
		return fmt.Errorf("recreate dst: %w", err)
	}
	if err := copyTree(src, dst); err != nil {
		return fmt.Errorf("copy to embed/public: %w", err)
	}
	log.Printf("Successfully exported ui dist to %s (%s)", dst, mode)
	return nil
}

func readPNPMVersion(packageJSON string) string {
	data, err := os.ReadFile(packageJSON)
	if err != nil {
		return ""
	}
	var payload struct {
		PackageManager string `json:"packageManager"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return ""
	}
	return strings.TrimPrefix(payload.PackageManager, "pnpm@")
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("go.mod not found")
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func runCmd(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func recreate(dir string) error {
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if entry.Name() == ".keep" || entry.Name() == "placeholder.txt" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(dir, entry.Name())); err != nil {
			return err
		}
	}
	return os.MkdirAll(dir, 0o755)
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		return err
	})
}
