package main

import (
  "embed"
)

//go:embed "Bee Dotty.ttf"
var embeddedFont []byte

//go:embed texts
var embeddedTexts embed.FS
