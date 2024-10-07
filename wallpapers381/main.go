package main

import (
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
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

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	window.SetKeyCallback(keyCallback)
	// respond to mouse movement
	window.SetCursorPosCallback(cursorPosCB)

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

func mouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	wWidth, wHeight := window.GetSize()

	// var widgetRS g143.Rect
	var widgetCode int

	for code, RS := range objCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	switch widgetCode {
	case PrevButton:
		if lineNo != 1 {
			lineNo = lineNo - 1
		}

		drawMainWindow(window, lineNo)

	case NextButton:
		lineNo = libw381.GetNextTextAddr(1)

		drawMainWindow(window, lineNo)

	case OurSite:

		if runtime.GOOS == "windows" {
			exec.Command("cmd", "/C", "start", "https://sae.ng").Run()
		} else if runtime.GOOS == "linux" {
			exec.Command("xdg-open", "https://sae.ng").Run()
		}

	case SetupInstrsButton:
		tmpFrame = currentWindowFrame
		dialogOpened = true
		drawSetupInstr(window, currentWindowFrame)

	case DialogCloseButton:
		if tmpFrame != nil {
			dialogOpened = false
			// send the frame to glfw window
			windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
			g143.DrawImage(wWidth, wHeight, tmpFrame, windowRS)
			window.SwapBuffers()

			currentWindowFrame = tmpFrame
			tmpFrame = nil
		}

	}
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	rootPath, _ := libw381.GetGUIPath()
	wWidth, wHeight := window.GetSize()

	setupInstrBtnRS := objCoords[SetupInstrsButton]
	entryRect := objCoords[WallpaperNumberEntry]

	// enforce number types
	if isKeyNumeric(key) {
		enteredText += glfw.GetKeyName(key, scancode)
	} else if key == glfw.KeyBackspace && len(enteredText) != 0 {
		enteredText = enteredText[:len(enteredText)-1]
	}

	enteredTextInt, _ := strconv.Atoi(enteredText)
	theCtx := Continue2dCtx(currentWindowFrame)
	theCtx.drawInput(WallpaperNumberEntry, entryRect.OriginX, 10, enteredTextInt)

	if key == glfw.KeyEnter {
		// check validity of entered number
		tmpAllTexts := strings.TrimSpace(string(libw381.EmbeddedTexts))
		numberOfTexts := len(strings.Split(tmpAllTexts, "\n"))

		tmp, err := strconv.Atoi(enteredText)
		if err != nil {
			return
		}
		if tmp > numberOfTexts {
			return
		}

		lineNo = tmp
		// update the image
		wimg := libw381.MakeAWallpaper(lineNo)
		w381OriginY := (setupInstrBtnRS.Height + 40)
		w381Width := wWidth - 20
		w381Height := wHeight - (w381OriginY)

		wimg = imaging.Fit(wimg, w381Width, w381Height, imaging.Lanczos)
		theCtx.ggCtx.DrawImage(wimg, 10, w381OriginY)

		os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte(strconv.Itoa(lineNo)), 0777)
	}

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}

func isKeyNumeric(key glfw.Key) bool {
	numKeys := []glfw.Key{glfw.Key0, glfw.Key1, glfw.Key2, glfw.Key3, glfw.Key4,
		glfw.Key5, glfw.Key6, glfw.Key7, glfw.Key8, glfw.Key9}

	for _, numKey := range numKeys {
		if key == numKey {
			return true
		}
	}

	return false
}

func cursorPosCB(window *glfw.Window, xpos, ypos float64) {
	if runtime.GOOS == "linux" {
		// linux fires too many events
		cursorEventsCount += 1
		if cursorEventsCount != 10 {
			return
		} else {
			cursorEventsCount = 0
		}
	}

	wWidth, wHeight := window.GetSize()

	var widgetRS g143.Rect
	var widgetCode int

	xPosInt := int(xpos)
	yPosInt := int(ypos)
	for code, RS := range objCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 || (widgetCode == DialogCloseButton && !dialogOpened) {
		// send the last drawn frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentWindowFrame, windowRS)
		window.SwapBuffers()
		return
	}

	rectA := image.Rect(widgetRS.OriginX, widgetRS.OriginY,
		widgetRS.OriginX+widgetRS.Width,
		widgetRS.OriginY+widgetRS.Height)

	pieceOfCurrentFrame := imaging.Crop(currentWindowFrame, rectA)
	invertedPiece := imaging.Invert(pieceOfCurrentFrame)

	ggCtx := gg.NewContextForImage(currentWindowFrame)
	ggCtx.DrawImage(invertedPiece, widgetRS.OriginX, widgetRS.OriginY)

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()
}
