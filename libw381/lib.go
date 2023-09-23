package libw381

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

const (
	DPI     = 72.0
	SIZE    = 90.0
	SPACING = 1.1
	MSIZE   = 45.0
)

func GetOutputTxt(lineNo int) string {
	tmpAllTexts := strings.TrimSpace(string(EmbeddedTexts))
	allTextsSlice := strings.Split(tmpAllTexts, "\n")
	return allTextsSlice[lineNo-1]
}

func MakeAWallpaper(lineNo int) image.Image {
	wWidth, wHeight := 1366, 768
	rgba := image.NewRGBA(image.Rect(0, 0, wWidth, wHeight))

	fontParsed, err := freetype.ParseFont(FontBytes)
	if err != nil {
		panic(err)
	}

	texts := GetOutputTxt(lineNo)

	newColor, _ := colorful.Hex("#3C2205")
	// newColor2, _ := colorful.Hex("#F2A550")

	fg := image.NewUniform(newColor)
	bg := image.NewUniform(color.White)

	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(DPI)
	c.SetFont(fontParsed)
	c.SetFontSize(SIZE)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	// Draw the text.
	pt := freetype.Pt(wWidth-200, wHeight-100)
	_, err = c.DrawString(strconv.Itoa(lineNo), pt)

	sizeW, sizeH := 60, 90
	currentX, currentY := 70, 50
	for _, cha := range texts {
		chaStr := strings.ToLower(string(cha))

		if chaStr != " " {
			if chaStr == "." {
				chaStr = "dot"
			} else if chaStr == "," {
				chaStr = "comma"
			} else if chaStr == "'" {
				chaStr = "apos"
			}

			rawCha, _ := Letters.ReadFile("letters/" + chaStr + ".png")

			chaImg, _, err := image.Decode(bytes.NewReader(rawCha))
			if err != nil {
				panic(err)
			}

			if chaStr == "dot" || chaStr == "comma" || chaStr == "apos" || chaStr == " " {
				chaImg = imaging.Fit(chaImg, sizeW/2, sizeH, imaging.Lanczos)
			} else {
				chaImg = imaging.Fit(chaImg, sizeW, sizeH, imaging.Lanczos)
			}

			draw.Draw(rgba, image.Rect(currentX, currentY, currentX+sizeW, currentY+sizeH), chaImg,
				image.Point{}, draw.Over)
		}

		if chaStr == "dot" || chaStr == "comma" || chaStr == "apos" || chaStr == " " {
			newX := currentX + 50
			if newX > (wWidth - 50) {
				currentY += sizeH
				currentX = 70
			} else {
				currentX += 50
			}

		} else {
			newX := currentX + sizeW
			if newX > (wWidth - sizeW - 20) {
				currentY += sizeH + 10
				currentX = 70
			} else {
				currentX += sizeW
			}

		}
	}

	return rgba
}

func GetDaemonPath() (string, error) {
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

func GetGUIPath() (string, error) {
	dd := os.Getenv("SNAP_USER_COMMON")
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

func GetNextTextAddr(method int) int {
	usePath := ""
	if method == 1 {
		tmp, _ := GetGUIPath()
		usePath = tmp
	} else {
		tmp, _ := GetDaemonPath()
		usePath = tmp
	}

	if DoesPathExists(filepath.Join(usePath, "last_text.txt")) {
		rawLastText, _ := os.ReadFile(filepath.Join(usePath, "last_text.txt"))
		number, err := strconv.Atoi(strings.TrimSpace(string(rawLastText)))
		if err != nil {
			panic(err)
		}
		toReturnNumber := number + 1
		tmpAllTexts := strings.TrimSpace(string(EmbeddedTexts))

		if toReturnNumber > len(strings.Split(tmpAllTexts, "\n")) {
			toReturnNumber = 1
		}
		os.WriteFile(filepath.Join(usePath, "last_text.txt"), []byte(strconv.Itoa(toReturnNumber)), 0777)
		return toReturnNumber
	} else {
		os.WriteFile(filepath.Join(usePath, "last_text.txt"), []byte("1"), 0777)
		return 1
	}

}
