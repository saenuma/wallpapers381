package libw381

import (
	"fmt"
	"image"
	"image/draw"
	"strings"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	DPI     = 72.0
	SIZE    = 90.0
	SPACING = 1.1
	MSIZE   = 45.0
)

func WordWrap(text string, writeWidth int, fontDrawer *font.Drawer) []string {
	widthFixed := fixed.I(writeWidth)

	strs := strings.Fields(text)
	outStrs := make([]string, 0)
	var tmpStr string
	for i, oneStr := range strs {
		var aStr string
		if i == 0 {
			aStr = oneStr
		} else {
			aStr += " " + oneStr
		}

		tmpStr += aStr
		if fontDrawer.MeasureString(tmpStr) >= widthFixed {
			outStr := tmpStr[:len(tmpStr)-len(aStr)]
			tmpStr = oneStr
			outStrs = append(outStrs, outStr)
		}
	}
	outStrs = append(outStrs, tmpStr)

	return outStrs
}

func GetOutputTxt(lineNo int) string {
	tmpAllTexts := strings.TrimSpace(string(EmbeddedTexts))
	allTextsSlice := strings.Split(tmpAllTexts, "\n")
	return fmt.Sprintf("%d. %s", lineNo, allTextsSlice[lineNo-1])
}

func MakeAWallpaper(lineNo int) image.Image {
	rgba := image.NewRGBA(image.Rect(0, 0, 1366, 768))

	fontParsed, err := freetype.ParseFont(FontBytes)
	if err != nil {
		panic(err)
	}

	fontDrawer := &font.Drawer{
		Dst: rgba,
		Src: image.Black,
		Face: truetype.NewFace(fontParsed, &truetype.Options{
			Size:    SIZE,
			DPI:     DPI,
			Hinting: font.HintingNone,
		}),
	}

	texts := WordWrap(GetOutputTxt(lineNo), 1366-130, fontDrawer)

	newColor, _ := colorful.Hex("#3C2205")
	newColor2, _ := colorful.Hex("#F2A550")

	fg := image.NewUniform(newColor)
	bg := image.NewUniform(newColor2)

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
	pt := freetype.Pt(80, 50+int(c.PointToFixed(SIZE)>>6))
	for _, s := range texts {
		_, err = c.DrawString(s, pt)
		if err != nil {
			panic(err)
		}
		pt.Y += c.PointToFixed(SIZE * SPACING)
	}

	return rgba
}
