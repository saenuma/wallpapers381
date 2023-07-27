package main

import (
	"bytes"
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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
	"github.com/saenuma/wallpapers381/libw381"
)

func main() {
	os.Setenv("FYNE_THEME", "light")

	rootPath, _ := libw381.GetGUIPath()

	myApp := app.New()
	myWindow := myApp.NewWindow("Wallpapers381 Gallery")

	// update slideshow store
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

	// gallery tab begin
	lineNo := libw381.GetNextTextAddr(1)
	wimg := libw381.MakeAWallpaper(lineNo)

	w381Img := canvas.NewImageFromImage(wimg)
	w381Img.FillMode = canvas.ImageFillOriginal

	imageContainer := container.NewCenter(container.NewPadded(w381Img))

	jumpEntry := widget.NewEntry()
	jumpEntry.SetText(strconv.Itoa(lineNo))
	jumpEntry.OnSubmitted = func(s string) {
		lineNo, err := strconv.Atoi(s)
		if err != nil {
			return
		}
		if lineNo > numberOfTexts {
			return
		}

		wimg = libw381.MakeAWallpaper(lineNo)
		w381Img = canvas.NewImageFromImage(wimg)
		w381Img.FillMode = canvas.ImageFillOriginal
		imageContainer.RemoveAll()
		imageContainer.Add(container.NewPadded(w381Img))
		os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte(strconv.Itoa(lineNo)), 0777)
	}

	nextBtn := widget.NewButton("next", func() {
		lineNo = libw381.GetNextTextAddr(1)
		wimg = libw381.MakeAWallpaper(lineNo)
		w381Img = canvas.NewImageFromImage(wimg)
		w381Img.FillMode = canvas.ImageFillOriginal
		imageContainer.RemoveAll()
		imageContainer.Add(container.NewPadded(w381Img))
		jumpEntry.SetText(strconv.Itoa(lineNo))
	})
	nextBtn.Importance = widget.HighImportance

	prevBtn := widget.NewButton("previous", func() {
		if lineNo != 1 {
			lineNo = lineNo - 1
		}

		wimg = libw381.MakeAWallpaper(lineNo)
		w381Img = canvas.NewImageFromImage(wimg)
		w381Img.FillMode = canvas.ImageFillOriginal
		imageContainer.RemoveAll()
		imageContainer.Add(container.NewPadded(w381Img))
		os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte(strconv.Itoa(lineNo)), 0777)
		jumpEntry.SetText(strconv.Itoa(lineNo))
	})

	setupInstructionsBtn := widget.NewButton("Setup Instructions", func() {

		var setupLabel *widget.RichText

		if runtime.GOOS == "windows" {
			hd, _ := os.UserHomeDir()
			path := filepath.Join(hd, "Wallpapers381")
			setupLabel = widget.NewRichTextFromMarkdown(fmt.Sprintf(`# Setup Instructions
		
1. Launch the App (needed to update the wallpapers store)
2. Open Settings.
3. Click **Personalisation** on the left and then click background
4. Set the first select to **Slideshow**
5. Click **Browse** and navigate to **%s** 
6. Repeat this instructions after update.
			`, path))

		} else {
			setupLabel = widget.NewRichTextFromMarkdown(`# Setup Instructions
1.	Launch the terminal

2.	Run the program **wallpapers381.switch**

			`)
		}

		innerBox := container.New(&FillSpace{}, container.NewMax(setupLabel))
		dialog.ShowCustom("Sample Lyrics File", "Close", innerBox, myWindow)
	})

	saeBtn := widget.NewButton("sae.ng", func() {
		if runtime.GOOS == "windows" {
			exec.Command("cmd", "/C", "start", "https://sae.ng").Run()
		} else if runtime.GOOS == "linux" {
			exec.Command("xdg-open", "https://sae.ng").Run()
		}
	})

	aboutBtn := widget.NewButton("About Us", func() {
		img, _, err := image.Decode(bytes.NewReader(SaeLogoBytes))
		if err != nil {
			panic(err)
		}
		logoImage := canvas.NewImageFromImage(img)
		logoImage.FillMode = canvas.ImageFillOriginal

		boxes := container.NewVBox(
			container.NewCenter(logoImage),
			widget.NewLabelWithStyle("Brought to You with Love by", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			saeBtn,
		)
		dialog.ShowCustom("About Wallpapers381", "Close", boxes, myWindow)
	})

	topBar := container.New(&halfes{}, prevBtn, nextBtn, jumpEntry, widget.NewSeparator(), setupInstructionsBtn, aboutBtn)
	galleryContainer := container.NewVBox(topBar, imageContainer)

	myWindow.SetContent(galleryContainer)

	myWindow.Resize(fyne.NewSize(1100, 600))
	myWindow.ShowAndRun()
}
