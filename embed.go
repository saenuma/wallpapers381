package main

import (
  "embed"
)

//go:embed fonts
var embeddedFonts embed.FS

//go:embed texts
var embeddedTexts embed.FS
