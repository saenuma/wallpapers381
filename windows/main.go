package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/saenuma/wallpapers381/libw381"
)

func main() {
	rootPath, _ := libw381.GetRootPath()

	myApp := app.New()
	myWindow := myApp.NewWindow("Wallpapers381 Gallery")

	lineNo := libw381.GetNextTextAddr()
	img := libw381.MakeAWallpaper(lineNo)

	w381Img := canvas.NewImageFromImage(img)
	w381Img.FillMode = canvas.ImageFillOriginal

	imageContainer := container.NewCenter(w381Img)

	tmpAllTexts := strings.TrimSpace(string(libw381.EmbeddedTexts))
	numberOfTexts := len(strings.Split(tmpAllTexts, "\n"))

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

		img = libw381.MakeAWallpaper(lineNo)
		w381Img = canvas.NewImageFromImage(img)
		w381Img.FillMode = canvas.ImageFillOriginal
		imageContainer.RemoveAll()
		imageContainer.Add(w381Img)
		os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte(strconv.Itoa(lineNo)), 0777)
	}

	nextBtn := widget.NewButton("next", func() {
		lineNo = libw381.GetNextTextAddr()
		img = libw381.MakeAWallpaper(lineNo)
		w381Img = canvas.NewImageFromImage(img)
		w381Img.FillMode = canvas.ImageFillOriginal
		imageContainer.RemoveAll()
		imageContainer.Add(w381Img)
		jumpEntry.SetText(strconv.Itoa(lineNo))
	})
	nextBtn.Importance = widget.HighImportance

	prevBtn := widget.NewButton("previous", func() {
		if lineNo != 1 {
			lineNo = lineNo - 1
		}

		img = libw381.MakeAWallpaper(lineNo)
		w381Img = canvas.NewImageFromImage(img)
		w381Img.FillMode = canvas.ImageFillOriginal
		imageContainer.RemoveAll()
		imageContainer.Add(w381Img)
		os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte(strconv.Itoa(lineNo)), 0777)
		jumpEntry.SetText(strconv.Itoa(lineNo))
	})

	bottomBar := container.New(&halfes{}, prevBtn, nextBtn, jumpEntry)
	galleryContainer := container.NewVBox(imageContainer, bottomBar)
	tabs := container.NewAppTabs(
		container.NewTabItem("Gallery", galleryContainer),
		container.NewTabItem("Setup Instructions", widget.NewLabel("World!")),
		container.NewTabItem("About Wallpapers381", widget.NewLabel("About")),
	)

	//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)

	myWindow.Resize(fyne.NewSize(1200, 700))
	myWindow.ShowAndRun()
}
