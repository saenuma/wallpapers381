package libw381

import (
	"embed"
)

//go:embed texts
var EmbeddedTexts embed.FS

//go:embed Barriecito-Regular.ttf
var FontBytes []byte
