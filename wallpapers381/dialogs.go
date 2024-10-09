package main

import (
	"image"
	"strings"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func drawSetupInstr(window *glfw.Window, currentFrame image.Image) {

	setupInstrStr := `
1.  Launch the terminal
2.  Run the program wallpapers381.switch
`

	// background image
	img := imaging.AdjustBrightness(currentFrame, -40)
	theCtx := Continue2dCtx(img)

	// dialog rectangle
	dialogWidth := 700
	dialogHeight := 600

	dialogOriginX := (theCtx.WindowWidth - dialogWidth) / 2
	dialogOriginY := (theCtx.WindowHeight - dialogHeight) / 2

	theCtx.ggCtx.SetHexColor("#fff")
	theCtx.ggCtx.DrawRectangle(float64(dialogOriginX), float64(dialogOriginY), float64(dialogWidth), float64(dialogHeight))
	theCtx.ggCtx.Fill()

	// load font
	fontPath := getDefaultFontPath()
	theCtx.ggCtx.LoadFontFace(fontPath, 40)

	// header message
	theCtx.ggCtx.SetHexColor("#444")
	h1 := "Setup Instructions"
	h1Width, h1Height := theCtx.ggCtx.MeasureString(h1)
	h1OriginX := (theCtx.WindowWidth - int(h1Width)) / 2
	theCtx.ggCtx.DrawString(h1, float64(h1OriginX), float64(dialogOriginY)+35+20)

	// message
	theCtx.ggCtx.LoadFontFace(fontPath, 20)

	for i, piece := range strings.Split(setupInstrStr, "\n") {
		pieceOriginY := dialogOriginY + 35 + int(h1Height) + (i+1)*20
		theCtx.ggCtx.DrawString(piece, float64(dialogOriginX)+20, float64(pieceOriginY))
	}

	// close button
	closeStr := "Close"
	closeStrWidth, _ := theCtx.ggCtx.MeasureString(closeStr)
	closeBtnOriginX := (theCtx.WindowWidth - int(closeStrWidth+50)) / 2
	closeBtnOriginY := dialogOriginY + dialogHeight - 50
	theCtx.drawButtonA(DialogCloseButton, closeBtnOriginX, closeBtnOriginY, closeStr, "#444", "#909BD0")

	// send the frame to glfw window
	windowRS := g143.Rect{Width: theCtx.WindowWidth, Height: theCtx.WindowHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(theCtx.WindowWidth, theCtx.WindowHeight, theCtx.ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}
