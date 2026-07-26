package static

import "embed"

//go:embed *.jpg *.css
var Assets embed.FS
