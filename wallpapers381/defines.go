package main

import (
	"image"

	g143 "github.com/bankole7782/graphics143"
)

const (
	fontSize = 20

	NextButton           = 101
	PrevButton           = 102
	WallpaperNumberEntry = 103
	SetupInstrsButton    = 104
	ColorsButton         = 105
	OurSite              = 106

	DialogCloseButton = 201
)

var (
	objCoords          map[int]g143.Rect
	currentWindowFrame image.Image
	lineNo             int
	enteredText        string
	tmpFrame           image.Image
	cursorEventsCount  = 0
	dialogOpened       bool
)
