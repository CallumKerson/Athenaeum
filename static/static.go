package static

import "embed"

//go:embed *.jpg *.css *.svg *.woff2 *.txt
var Assets embed.FS
