package main

import (
  "embed"
)

//go:embed "JandaQuirkygirl.ttf"
var embeddedFont []byte

//go:embed texts
var embeddedTexts embed.FS
