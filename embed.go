package main

import (
  "embed"
)

//go:embed AmaticSC-Bold.ttf
var embeddedFont []byte

//go:embed texts
var embeddedTexts embed.FS
