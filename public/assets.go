package assets

import "embed"

var embedAssets embed.FS

func GetAssets() embed.FS {
	return embedAssets
}

