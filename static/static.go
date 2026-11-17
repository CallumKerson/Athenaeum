package static

import "embed"

//go:embed *.jpg *.css *.svg *.woff2 *.txt *.js
var Assets embed.FS
