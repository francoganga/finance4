package templates

import "embed"

//go:embed all:*
var embedTemplates embed.FS

func GetTemplates() embed.FS {
	return embedTemplates
}

