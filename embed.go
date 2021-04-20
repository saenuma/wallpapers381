package main

import (
  "embed"
)

//go:embed "Drifttype.ttf"
var embeddedFont []byte

//go:embed texts
var embeddedTexts embed.FS
