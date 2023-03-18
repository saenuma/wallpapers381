package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Wallpapers381 Gallery")

	tabs := container.NewAppTabs(
		container.NewTabItem("Gallery", widget.NewLabel("Hello")),
		container.NewTabItem("Setup Instructions", widget.NewLabel("World!")),
	)

	//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)

	myWindow.Resize(fyne.NewSize(1200, 700))
	myWindow.ShowAndRun()
}
