package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	DialogCloseButton = 201
)

func drawSetupInstr(window *glfw.Window, currentFrame image.Image) {

	var setupInstrStr string
	if runtime.GOOS == "windows" {
		hd, _ := os.UserHomeDir()
		path := filepath.Join(hd, "Wallpapers381")
		setupInstrStr = fmt.Sprintf(`
1. Launch the App (needed to update the wallpapers store)
2. Open Settings.
3. Click Personalisation on the left and then click background
4. Set the first select to Slideshow
5. Click Browse and navigate to %s 
6. Repeat this instructions after update.
`, path)

	} else {
		setupInstrStr = `
1.  Launch the terminal
2.  Run the program wallpapers381.switch
`
	}

	wWidth, wHeight := window.GetSize()

	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background image
	img := imaging.AdjustBrightness(currentFrame, -40)
	ggCtx.DrawImage(img, 0, 0)

	// load font
	fontPath := getDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 40)
	if err != nil {
		panic(err)
	}

	// dialog rectangle
	dialogWidth := 700
	dialogHeight := 600

	dialogOriginX := (wWidth - dialogWidth) / 2
	dialogOriginY := (wHeight - dialogHeight) / 2

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawRectangle(float64(dialogOriginX), float64(dialogOriginY), float64(dialogWidth), float64(dialogHeight))
	ggCtx.Fill()

	// header message
	ggCtx.SetHexColor("#444")
	h1 := "Setup Instructions"
	h1Width, h1Height := ggCtx.MeasureString(h1)
	h1OriginX := (wWidth - int(h1Width)) / 2
	ggCtx.DrawString(h1, float64(h1OriginX), float64(dialogOriginY)+35+20)

	// message
	err = ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	for i, piece := range strings.Split(setupInstrStr, "\n") {
		pieceOriginY := dialogOriginY + 35 + int(h1Height) + (i+1)*20
		ggCtx.DrawString(piece, float64(dialogOriginX)+20, float64(pieceOriginY))
	}

	// close button
	closeStr := "Close"
	closeStrWidth, closeStrHeight := ggCtx.MeasureString(closeStr)
	ggCtx.SetHexColor("#909BD0")
	closeBtnOriginX := (wWidth - int(closeStrWidth+50)) / 2
	ggCtx.DrawRoundedRectangle(float64(closeBtnOriginX), float64(dialogOriginY+dialogHeight-50), closeStrWidth+50,
		closeStrHeight+25, (closeStrHeight+25)/2)
	ggCtx.Fill()

	closeBtnRS := g143.RectSpecs{Width: int(closeStrWidth) + 50, Height: int(closeStrHeight) + 25,
		OriginX: closeBtnOriginX, OriginY: dialogOriginY + dialogHeight - 50}
	objCoords[DialogCloseButton] = closeBtnRS

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(closeStr, float64(closeBtnOriginX+25), float64(dialogOriginY+dialogHeight-25))

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = ggCtx.Image()
}
