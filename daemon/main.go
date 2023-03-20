package main

import (
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/saenuma/wallpapers381/libw381"
)

func main() {

	for {
		lineNo := libw381.GetNextTextAddr(2)
		img := libw381.MakeAWallpaper(lineNo)
		outPath := getOutputPath()

		imaging.Save(img, outPath)

		// sleep for 30 seconds
		time.Sleep(30 * 60 * time.Second)
		// time.Sleep(1 * 60 * time.Second)

	}

}

func getOutputPath() string {
	rootPath, _ := libw381.GetDaemonPath()
	return filepath.Join(rootPath, "wallpaper.png")
}
