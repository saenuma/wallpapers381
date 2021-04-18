package main

import (
  "embed"
)

//go:embed "Melted Monster.ttf"
var embeddedFont []byte

//go:embed texts
var embeddedTexts embed.FS
