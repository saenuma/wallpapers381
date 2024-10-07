package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kovidgoyal/imaging"
	"github.com/saenuma/wallpapers381/libw381"
)

func main() {
	runtime.LockOSThread()

	objCoords = make(map[int]g143.Rect)

	window := g143.NewWindow(1200, 800, "Wallpapers381 Gallery", false)

	lineNo = libw381.GetNextTextAddr(1)
	drawMainWindow(window, lineNo)

	// // respond to the mouse
	// window.SetMouseButtonCallback(mouseBtnCallback)
	// window.SetKeyCallback(keyCallback)
	// // respond to mouse movement
	// window.SetCursorPosCallback(cursorPosCB)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(24) - time.Since(t))
	}

}

func getDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "k117_font.ttf")
	os.WriteFile(fontPath, DefaultFont, 0777)
	return fontPath
}

func drawMainWindow(window *glfw.Window, lineNo int) {
	wWidth, wHeight := window.GetSize()
	theCtx := New2dCtx(wWidth, wHeight)

	prevBtnRect := theCtx.drawButtonA(PrevButton, 200, 10, "Previous", "#444", "#D09090")
	nextBtnOriginX, nextBtnOriginY := nextHorizontalCoords(prevBtnRect, 30)
	nextBtnRect := theCtx.drawButtonA(NextButton, nextBtnOriginX, nextBtnOriginY, "Next", "#444", "#90D092")

	numStr := "Wallpaper Number:"
	numStrWidth, _ := theCtx.ggCtx.MeasureString(numStr)
	numStrOriginX, _ := nextHorizontalCoords(nextBtnRect, 30)
	theCtx.ggCtx.DrawString(numStr, float64(numStrOriginX), 35)

	entryX := numStrWidth + float64(numStrOriginX) + 30

	lineNoStr := strconv.Itoa(lineNo)
	enteredText = lineNoStr

	entryRect := theCtx.drawInput(WallpaperNumberEntry, int(entryX), 10, lineNo)
	setupBtnX, _ := nextHorizontalCoords(entryRect, 30)
	setupBtnRect := theCtx.drawButtonA(SetupInstrsButton, setupBtnX, nextBtnOriginY, "Setup Instructions", "#444", "#D090CB")

	// display current wallpaper
	wimg := libw381.MakeAWallpaper(lineNo)

	w381OriginY := (setupBtnRect.Height + 40)
	w381Width := wWidth - 20
	w381Height := wHeight - (w381OriginY)

	wimg = imaging.Fit(wimg, w381Width, w381Height, imaging.Lanczos)
	theCtx.ggCtx.DrawImage(wimg, 10, w381OriginY)

	theCtx.drawButtonB(OurSite, "https://sae.ng/", "#9C5858")

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}
