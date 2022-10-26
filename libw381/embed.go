package libw381

import (
	_ "embed"
)

//go:embed msgs.txt
var EmbeddedTexts []byte

//go:embed Barriecito-Regular.ttf
var FontBytes []byte
