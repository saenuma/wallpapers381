package main

import (
  "math/rand"
  "time"
  "fmt"
  "image"
  "path/filepath"
  "image/color"
	"image/draw"
	"image/png"
	"os"
  "github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
  "golang.org/x/image/math/fixed"
  "golang.org/x/image/font"
  "strings"
  "bufio"
  "github.com/go-playground/colors"
  "github.com/kbinani/screenshot"
)


const (
  DPI = 72.0
  SIZE = 90.0
  SPACING = 1.1
)


func main() {
  fontParsed, err := freetype.ParseFont(embeddedFont)
  if err != nil {
    panic(err)
  }
  hex, err := colors.ParseHEX("#C9B466")
  nCR := hex.ToRGBA()
  newColor := color.RGBA{uint8(nCR.R), uint8(nCR.G), uint8(nCR.B), uint8(nCR.A)}
  bg, fg := image.Black, image.NewUniform(newColor)


  screenBounds := screenshot.GetDisplayBounds(0)

  rgba := image.NewRGBA(image.Rect(0, 0, screenBounds.Dx(), screenBounds.Dy()))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(DPI)
	c.SetFont(fontParsed)
	c.SetFontSize(SIZE)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
  c.SetHinting(font.HintingNone)

  fontDrawer := &font.Drawer{
    Dst: rgba,
    Src: image.Black,
    Face: truetype.NewFace(fontParsed, &truetype.Options{
      Size: SIZE,
      DPI: DPI,
      Hinting: font.HintingNone,
    }),
  }

  textBytes, err := embeddedTexts.ReadFile(randTextFile())
  if err != nil {
    panic(err)
  }
  texts := wordWrap(string(textBytes), screenBounds.Dx() - 200, fontDrawer)

  // Draw the text.
	pt := freetype.Pt(100, 100+int(c.PointToFixed(SIZE)>>6))
	for _, s := range texts {
		_, err = c.DrawString(s, pt)
		if err != nil {
			panic(err)
		}
		pt.Y += c.PointToFixed(SIZE * SPACING)
	}

  // Save that RGBA image to disk.
	outFile, err := os.Create("/tmp/out.png")
	if err != nil {
    panic(err)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
    panic(err)
	}
	err = b.Flush()
	if err != nil {
    panic(err)
	}
}


func randTextFile() string {
  dirFIs, err := embeddedTexts.ReadDir("texts")
  if err != nil {
    panic(err)
  }
  fonts := make([]string, 0)
  for _, dirFI := range dirFIs {
    f := filepath.Join("texts", dirFI.Name())
    fonts = append(fonts, f)
  }

  var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
  return fonts[seededRand.Intn(len(fonts))]
}


func wordWrap(text string, writeWidth int, fontDrawer *font.Drawer) []string {
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
      outStr := tmpStr[ : len(tmpStr) - len(aStr) ]
      tmpStr = oneStr
      outStrs = append(outStrs, outStr)
    }
  }
  outStrs = append(outStrs, tmpStr)

  return outStrs
}
