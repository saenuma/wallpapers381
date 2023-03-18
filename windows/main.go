package main

import (
	"os"
	"path/filepath"
	"strconv"

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

	nextBtn := widget.NewButton("next", func() {
		lineNo = libw381.GetNextTextAddr()
		img = libw381.MakeAWallpaper(lineNo)
		w381Img = canvas.NewImageFromImage(img)
		w381Img.FillMode = canvas.ImageFillOriginal
		imageContainer.RemoveAll()
		imageContainer.Add(w381Img)
	})
	nextBtn.Importance = widget.HighImportance

	prevBtn := widget.NewButton("previous", func() {
		if lineNo != 1 {
			lineNo = lineNo - 1
		}
		os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte(strconv.Itoa(lineNo)), 0777)

		img = libw381.MakeAWallpaper(lineNo)
		w381Img = canvas.NewImageFromImage(img)
		w381Img.FillMode = canvas.ImageFillOriginal
		imageContainer.RemoveAll()
		imageContainer.Add(w381Img)
	})

	bottomBar := container.New(&halfes{}, prevBtn, nextBtn)
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
