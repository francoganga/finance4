package assets

import "embed"

//go:embed *
var embedAssets embed.FS

func GetAssets() embed.FS {
	return embedAssets
}

