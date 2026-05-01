//go:build embed

package web

import (
	"embed"
	"io/fs"
)

//go:embed embed/public
var embeddedFS embed.FS

// FS exposes the embedded SPA assets for production builds.
func FS() fs.FS {
	sub, err := fs.Sub(embeddedFS, "embed/public")
	if err != nil {
		panic(err)
	}
	return sub
}
