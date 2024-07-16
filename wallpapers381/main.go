package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/wallpapers381/libw381"
)

const (
	fps                  = 24
	NextButton           = 101
	PrevButton           = 102
	WallpaperNumberEntry = 103
	SetupInstrsButton    = 104
	ColorsButton         = 105
	OurSite              = 106
)

var objCoords map[int]g143.RectSpecs
var currentWindowFrame image.Image
var lineNo int
var enteredText string
var tmpFrame image.Image
var cursorEventsCount = 0

func main() {
	rootPath, _ := libw381.GetGUIPath()
	// update slideshow store for windows
	tmpAllTexts := strings.TrimSpace(string(libw381.EmbeddedTexts))
	numberOfTexts := len(strings.Split(tmpAllTexts, "\n"))

	if runtime.GOOS == "windows" {
		numberOfCPUS := runtime.NumCPU()
		var wg sync.WaitGroup
		jobsPerThread := int(math.Floor(float64(numberOfTexts) / float64(numberOfCPUS)))

		installedVersion := ""
		rawVersion, err := os.ReadFile(filepath.Join(rootPath, "version.txt"))
		if err != nil {
			installedVersion = "undefined"
		}
		installedVersion = strings.TrimSpace(string(rawVersion))

		if W381_IMAGES_VERSION != installedVersion {
			hd, _ := os.UserHomeDir()
			if libw381.DoesPathExists(filepath.Join(hd, "Wallpapers381")) {
				os.RemoveAll(filepath.Join(hd, "Wallpapers381"))
			}
			os.MkdirAll(filepath.Join(hd, "Wallpapers381"), 0777)

			for threadIndex := 0; threadIndex < numberOfCPUS; threadIndex++ {
				wg.Add(1)
				startIndex := threadIndex * jobsPerThread
				endIndex := (threadIndex + 1) * jobsPerThread

				go func(startIndex, endIndex int, wg *sync.WaitGroup) {
					defer wg.Done()

					for index := startIndex; index < endIndex; index++ {
						if index == 0 {
							continue
						}

						img := libw381.MakeAWallpaper(index)
						imaging.Save(img, filepath.Join(hd, "Wallpapers381", fmt.Sprintf("%d.png", index)))
					}
				}(startIndex, endIndex, &wg)
			}
			wg.Wait()

			for index := (jobsPerThread * numberOfCPUS); index < numberOfTexts; index++ {
				img := libw381.MakeAWallpaper(index)
				imaging.Save(img, filepath.Join(hd, "Wallpapers381", fmt.Sprintf("%d.png", index)))
			}

			os.WriteFile(filepath.Join(rootPath, "version.txt"), []byte(W381_IMAGES_VERSION), 0777)
		}
	}

	runtime.LockOSThread()

	objCoords = make(map[int]g143.RectSpecs)

	window := g143.NewWindow(1200, 800, "Wallpapers381 Gallery", false)

	lineNo = libw381.GetNextTextAddr(1)
	allDraws(window, lineNo)

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	window.SetKeyCallback(keyCallback)
	// respond to mouse movement
	window.SetCursorPosCallback(cursorPosCB)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func allDraws(window *glfw.Window, lineNo int) {
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
	ggCtx.DrawRectangle(float64(beginXOffset), 10, prevStrW+50, prevStrH+25)
	ggCtx.Fill()

	prevBtnRS := g143.RectSpecs{Width: int(prevStrW) + 50, Height: int(prevStrH) + 25, OriginX: beginXOffset, OriginY: 10}
	objCoords[PrevButton] = prevBtnRS

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(prevStr, float64(beginXOffset)+25, 35)

	// next button
	ggCtx.SetHexColor("#90D092")
	nextStr := "Next"
	nextStrWidth, nextStrHeight := ggCtx.MeasureString(nextStr)
	nexBtnOriginX := prevBtnRS.OriginX + prevBtnRS.Width + 30
	ggCtx.DrawRectangle(float64(nexBtnOriginX), 10, nextStrWidth+50, nextStrHeight+25)
	ggCtx.Fill()

	nextBtnRS := g143.RectSpecs{Width: int(nextStrWidth) + 50, Height: int(nextStrHeight) + 25, OriginX: nexBtnOriginX,
		OriginY: 10}
	objCoords[NextButton] = nextBtnRS

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
	objCoords[WallpaperNumberEntry] = wNumEntryRS

	lineNoStr := strconv.Itoa(lineNo)
	enteredText = lineNoStr
	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(lineNoStr, float64(wNumEntryOriginX+15), 35)

	// setup instructions button
	ggCtx.SetHexColor("#D090CB")
	setupInstrStr := "Setup Instructions"
	setupInstrStrWidth, setupInstrStrHeight := ggCtx.MeasureString(setupInstrStr)
	setupInstrBtnOriginX := wNumEntryRS.OriginX + wNumEntryRS.Width + 30
	ggCtx.DrawRectangle(float64(setupInstrBtnOriginX), 10, setupInstrStrWidth+50,
		setupInstrStrHeight+25)
	ggCtx.Fill()

	setupInstrBtnRS := g143.RectSpecs{Width: int(setupInstrStrWidth) + 50, Height: int(setupInstrStrHeight) + 25,
		OriginX: setupInstrBtnOriginX, OriginY: 10}
	objCoords[SetupInstrsButton] = setupInstrBtnRS

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
	objCoords[OurSite] = fars

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

	// var widgetRS g143.RectSpecs
	var widgetCode int

	for code, RS := range objCoords {
		if g143.InRectSpecs(RS, xPosInt, yPosInt) {
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

		allDraws(window, lineNo)

	case NextButton:
		lineNo = libw381.GetNextTextAddr(1)

		allDraws(window, lineNo)

	case OurSite:

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

	rootPath, _ := libw381.GetGUIPath()
	wWidth, wHeight := window.GetSize()

	setupInstrBtnRS := objCoords[SetupInstrsButton]
	wNumEntryRS := objCoords[WallpaperNumberEntry]

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

	var widgetRS g143.RectSpecs
	var widgetCode int

	xPosInt := int(xpos)
	yPosInt := int(ypos)
	for code, RS := range objCoords {
		if g143.InRectSpecs(RS, xPosInt, yPosInt) {
			widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		// send the last drawn frame to glfw window
		windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
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
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()
}
