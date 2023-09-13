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
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/wallpapers381/libw381"
)

const (
	fps = 10
)

var objCoords map[g143.RectSpecs]any
var currentWindowFrame image.Image
var lineNo int
var wNumEntryActive bool
var enteredText string
var tmpFrame image.Image

// symbols types
type NextButton struct{}
type PrevButton struct{}
type WallpaperNumberEntry struct{}
type SetupInstrsButton struct{}
type ColorsButton struct{}
type OurSite struct{}

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
	beginXOffset := 200
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

	lineNo = libw381.GetNextTextAddr(1)
	lineNoStr := strconv.Itoa(lineNo)
	enteredText = lineNoStr
	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(lineNoStr, float64(wNumEntryOriginX+15), 35)

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

	// display current wallpaper
	wimg := libw381.MakeAWallpaper(lineNo)

	w381OriginY := (setupInstrBtnRS.Height + 40)
	w381Width := wWidth - 20
	w381Height := wHeight - (w381OriginY)

	wimg = imaging.Fit(wimg, w381Width, w381Height, imaging.Lanczos)
	ggCtx.DrawImage(wimg, 10, w381OriginY)

	// draw our site below
	ggCtx.SetHexColor("#9C5858")
	fromAddr := "sae.ng"
	fromAddrWidth, fromAddrHeight := ggCtx.MeasureString(fromAddr)
	fromAddrOriginX := (wWidth - int(fromAddrWidth)) / 2
	ggCtx.DrawString(fromAddr, float64(fromAddrOriginX), float64(wHeight-int(fromAddrHeight)))
	fars := g143.RectSpecs{OriginX: fromAddrOriginX, OriginY: wHeight - 40,
		Width: int(fromAddrWidth), Height: 40}
	objCoords[fars] = OurSite{}

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

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	wWidth, wHeight := window.GetSize()

	// var objRS g143.RectSpecs
	var obj any

	var setupInstrBtnRS, wNumEntryRS g143.RectSpecs

	for rs, anyObj := range objCoords {
		if g143.InRectSpecs(rs, xPosInt, yPosInt) {
			// objRS = rs
			obj = anyObj
		}

		// store setupInstrBtnRS
		if _, ok := anyObj.(SetupInstrsButton); ok {
			setupInstrBtnRS = rs
		}

		// store wNumEntryRS
		if _, ok := anyObj.(WallpaperNumberEntry); ok {
			wNumEntryRS = rs
		} else {
			wNumEntryActive = false
		}
	}

	rootPath, _ := libw381.GetGUIPath()

	if obj == nil {
		return
	}

	switch obj.(type) {
	case PrevButton:
		wNumEntryActive = false
		if lineNo != 1 {
			lineNo = lineNo - 1
		}

		ggCtx := gg.NewContextForImage(currentWindowFrame)

		// update the image
		wimg := libw381.MakeAWallpaper(lineNo)
		w381OriginY := (setupInstrBtnRS.Height + 40)
		w381Width := wWidth - 20
		w381Height := wHeight - (w381OriginY)

		wimg = imaging.Fit(wimg, w381Width, w381Height, imaging.Lanczos)
		ggCtx.DrawImage(wimg, 10, w381OriginY)

		// load font
		fontPath := getDefaultFontPath()
		err := ggCtx.LoadFontFace(fontPath, 20)
		if err != nil {
			panic(err)
		}

		// update the display of line number
		lineNoStr := strconv.Itoa(lineNo)
		enteredText = lineNoStr
		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(float64(wNumEntryRS.OriginX), 10,
			float64(wNumEntryRS.Width), float64(setupInstrBtnRS.Height-15))
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(lineNoStr, float64(wNumEntryRS.OriginX+15), 35)
		os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte(strconv.Itoa(lineNo)), 0777)

		// send the frame to glfw window
		windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
		window.SwapBuffers()

		// save the frame
		currentWindowFrame = ggCtx.Image()

	case NextButton:
		wNumEntryActive = false
		lineNo = libw381.GetNextTextAddr(1)

		ggCtx := gg.NewContextForImage(currentWindowFrame)

		// update the image
		wimg := libw381.MakeAWallpaper(lineNo)
		w381OriginY := (setupInstrBtnRS.Height + 40)
		w381Width := wWidth - 20
		w381Height := wHeight - (w381OriginY)

		wimg = imaging.Fit(wimg, w381Width, w381Height, imaging.Lanczos)
		ggCtx.DrawImage(wimg, 10, w381OriginY)

		// load font
		fontPath := getDefaultFontPath()
		err := ggCtx.LoadFontFace(fontPath, 20)
		if err != nil {
			panic(err)
		}

		// update the display of line number
		lineNoStr := strconv.Itoa(lineNo)
		enteredText = lineNoStr
		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(float64(wNumEntryRS.OriginX), 10,
			float64(wNumEntryRS.Width), float64(setupInstrBtnRS.Height-15))
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(lineNoStr, float64(wNumEntryRS.OriginX+15), 35)
		os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte(strconv.Itoa(lineNo)), 0777)

		// send the frame to glfw window
		windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
		window.SwapBuffers()

		// save the frame
		currentWindowFrame = ggCtx.Image()

	case WallpaperNumberEntry:
		wNumEntryActive = true

	case OurSite:
		wNumEntryActive = false

		if runtime.GOOS == "windows" {
			exec.Command("cmd", "/C", "start", "https://sae.ng").Run()
		} else if runtime.GOOS == "linux" {
			exec.Command("xdg-open", "https://sae.ng").Run()
		}

	case SetupInstrsButton:
		tmpFrame = currentWindowFrame

		drawSetupInstr(window, currentWindowFrame)

	case DialogCloseButton:
		if tmpFrame != nil {
			// send the frame to glfw window
			windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
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

	if !wNumEntryActive {
		return
	}

	rootPath, _ := libw381.GetGUIPath()
	wWidth, wHeight := window.GetSize()

	var wNumEntryRS, setupInstrBtnRS g143.RectSpecs
	for k, v := range objCoords {
		if _, ok := v.(WallpaperNumberEntry); ok {
			wNumEntryRS = k
		}
		if _, ok := v.(SetupInstrsButton); ok {
			setupInstrBtnRS = k
		}
	}

	// enforce number types
	if isKeyNumeric(key) {
		enteredText += glfw.GetKeyName(key, scancode)
	} else if key == glfw.KeyBackspace && len(enteredText) != 0 {
		enteredText = enteredText[:len(enteredText)-1]
	}

	fontPath := getDefaultFontPath()
	ggCtx := gg.NewContextForImage(currentWindowFrame)
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawRectangle(float64(wNumEntryRS.OriginX), 10,
		float64(wNumEntryRS.Width), float64(setupInstrBtnRS.Height-15))
	ggCtx.Fill()

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(enteredText, float64(wNumEntryRS.OriginX+15), 35)

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
		ggCtx.DrawImage(wimg, 10, w381OriginY)

		os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte(strconv.Itoa(lineNo)), 0777)
	}

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = ggCtx.Image()
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