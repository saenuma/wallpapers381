package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/saenuma/wallpapers381/libw381"
)

func main() {

	for {
		lineNo := getNextTextAddr()
		img := libw381.MakeAWallpaper(lineNo)
		outPath := getOutputPath()

		imaging.Save(img, outPath)

		// sleep for 30 seconds
		time.Sleep(30 * 60 * time.Second)
		// time.Sleep(1 * 60 * time.Second)

	}

}

func GetRootPath() (string, error) {
	dd := os.Getenv("SNAP_COMMON")
	if strings.HasPrefix(dd, "/var/snap/go") || dd == "" {
		hd, err := os.UserHomeDir()
		if err != nil {
			return "", errors.Wrap(err, "os error")
		}
		dd = filepath.Join(hd, "W381")
		os.MkdirAll(dd, 0777)
	}

	return dd, nil
}

func DoesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

func getNextTextAddr() int {
	rootPath, _ := GetRootPath()
	if DoesPathExists(filepath.Join(rootPath, "last_text.txt")) {
		rawLastText, _ := os.ReadFile(filepath.Join(rootPath, "last_text.txt"))
		number, err := strconv.Atoi(strings.TrimSpace(string(rawLastText)))
		if err != nil {
			panic(err)
		}
		toReturnNumber := number + 1
		tmpAllTexts := strings.TrimSpace(string(libw381.EmbeddedTexts))

		if toReturnNumber > len(strings.Split(tmpAllTexts, "\n")) {
			toReturnNumber = 1
		}
		os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte(strconv.Itoa(toReturnNumber)), 0777)
		return toReturnNumber
	} else {
		os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte("1"), 0777)
		return 1
	}

}

func getOutputPath() string {
	rootPath, _ := GetRootPath()
	return filepath.Join(rootPath, "wallpaper.png")
}
