package main

import (
	"image"
	"os"
	"path/filepath"
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	fps = 10
)

var objCoords map[g143.RectSpecs]any
var currentWindowFrame image.Image

// symbols types
type NextButton struct{}
type PrevButton struct{}
type WallpaperNumberEntry struct{}
type SetupInstrsButton struct{}
type ColorsButton struct{}

func main() {
	runtime.LockOSThread()

	objCoords = make(map[g143.RectSpecs]any)

	window := g143.NewWindow(1200, 800, "Wallpapers381 Gallery", false)
	allDraws(window)

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	window.SetKeyCallback(keyCallback)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func allDraws(window *glfw.Window) {
	wWidth, wHeight := window.GetSize()

	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := getDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	// previous button
	beginXOffset := 80
	ggCtx.SetHexColor("#D09090")
	prevStr := "Previous"
	prevStrW, prevStrH := ggCtx.MeasureString(prevStr)
	ggCtx.DrawRoundedRectangle(float64(beginXOffset), 10, prevStrW+50, prevStrH+25, (prevStrH+25)/2)
	ggCtx.Fill()

	prevBtnRS := g143.RectSpecs{Width: int(prevStrW) + 50, Height: int(prevStrH) + 25, OriginX: beginXOffset, OriginY: 10}
	objCoords[prevBtnRS] = PrevButton{}

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(prevStr, float64(beginXOffset)+25, 35)

	// next button
	ggCtx.SetHexColor("#90D092")
	nextStr := "Next"
	nextStrWidth, nextStrHeight := ggCtx.MeasureString(nextStr)
	nexBtnOriginX := prevBtnRS.OriginX + prevBtnRS.Width + 30
	ggCtx.DrawRoundedRectangle(float64(nexBtnOriginX), 10, nextStrWidth+50, nextStrHeight+25, (nextStrHeight+25)/2)
	ggCtx.Fill()

	nextBtnRS := g143.RectSpecs{Width: int(nextStrWidth) + 50, Height: int(nextStrHeight) + 25, OriginX: nexBtnOriginX,
		OriginY: 10}
	objCoords[nextBtnRS] = NextButton{}

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(nextStr, float64(nextBtnRS.OriginX)+25, 35)

	numStr := "Wallpaper Number:"
	numStrWidth, _ := ggCtx.MeasureString(numStr)
	ggCtx.DrawString(numStr, float64(nextBtnRS.OriginX+nextBtnRS.Width)+30, 35)

	// wallpaper number entry box
	ggCtx.SetHexColor("#909BD0")
	wNumEntryOriginX := nextBtnRS.OriginX + nextBtnRS.Width + 30 + int(numStrWidth) + 10
	ggCtx.DrawRectangle(float64(wNumEntryOriginX), nextStrHeight+30, 100, 3)
	ggCtx.Fill()

	wNumEntryRS := g143.RectSpecs{Width: 100, Height: int(nextStrHeight) + 30, OriginX: wNumEntryOriginX,
		OriginY: 10}
	objCoords[wNumEntryRS] = WallpaperNumberEntry{}

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString("1", float64(wNumEntryOriginX+15), 35)

	// setup instructions button
	ggCtx.SetHexColor("#D090CB")
	setupInstrStr := "Setup Instructions"
	setupInstrStrWidth, setupInstrStrHeight := ggCtx.MeasureString(setupInstrStr)
	setupInstrBtnOriginX := wNumEntryRS.OriginX + wNumEntryRS.Width + 30
	ggCtx.DrawRoundedRectangle(float64(setupInstrBtnOriginX), 10, setupInstrStrWidth+50,
		setupInstrStrHeight+25, (setupInstrStrHeight+25)/2)
	ggCtx.Fill()

	setupInstrBtnRS := g143.RectSpecs{Width: int(setupInstrStrWidth) + 50, Height: int(setupInstrStrHeight) + 25,
		OriginX: setupInstrBtnOriginX, OriginY: 10}
	objCoords[setupInstrBtnRS] = SetupInstrsButton{}

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(setupInstrStr, float64(setupInstrBtnOriginX+25), 35)

	// colors button
	ggCtx.SetHexColor("#D0BB90")
	colorsStr := "Setup Colors"
	colorsStrWidth, colorsStrHeight := ggCtx.MeasureString(colorsStr)
	colorsBtnOriginX := setupInstrBtnOriginX + setupInstrBtnRS.Width + 30
	ggCtx.DrawRoundedRectangle(float64(colorsBtnOriginX), 10, colorsStrWidth+50, colorsStrHeight+25,
		(colorsStrHeight+25)/2)
	ggCtx.Fill()

	colorsBtnRS := g143.RectSpecs{Width: int(colorsStrWidth) + 50, Height: int(colorsStrHeight) + 25,
		OriginX: colorsBtnOriginX, OriginY: 10}
	objCoords[colorsBtnRS] = ColorsButton{}

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(colorsStr, float64(colorsBtnOriginX+25), 35)

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = ggCtx.Image()
}

func getDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "k117_font.ttf")
	os.WriteFile(fontPath, DefaultFont, 0777)
	return fontPath
}

func mouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

}
