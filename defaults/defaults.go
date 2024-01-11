package defaults

import (
	"embed"
)

//go:embed *.html
var DefaultTemplates embed.FS
