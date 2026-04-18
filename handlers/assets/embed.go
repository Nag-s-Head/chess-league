package assets

import (
	"embed"

	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

//go:embed *.jpg

var fs embed.FS

var wideShot []byte = utils.GetRaw(fs, "wide_shot.jpg")
