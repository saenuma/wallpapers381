package libw381

import (
	"embed"
)

//go:embed msgs.txt
var EmbeddedTexts []byte

//go:embed Barriecito-Regular.ttf
var FontBytes []byte

//go:embed letters
var Letters embed.FS
