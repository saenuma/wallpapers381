package main

import (
  "math/rand"
  "time"
  // "fmt"
  "image"
  "path/filepath"
  // "image/color"
	"image/draw"
	"image/png"
	"os"
	"github.com/golang/freetype"
  "golang.org/x/image/font"
  "strings"
  "bufio"
)


const (
  DPI = 72.0
  SIZE = 80.0
  SPACING = 1.0
)


func main() {
  textBytes, err := embeddedTexts.ReadFile(randTextFile())
  if err != nil {
    panic(err)
  }
  text := strings.Split(string(textBytes), "\n")

  fontBytes, err := embeddedFonts.ReadFile(randFontFile())
  if err != nil {
    panic(err)
  }
  fontParsed, err := freetype.ParseFont(fontBytes)
  if err != nil {
    panic(err)
  }
  fg, bg := image.Black, image.White

  rgba := image.NewRGBA(image.Rect(0, 0, 1366, 768))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(DPI)
	c.SetFont(fontParsed)
	c.SetFontSize(SIZE)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
  c.SetHinting(font.HintingNone)

  // Draw the text.
	pt := freetype.Pt(50, 50+int(c.PointToFixed(SIZE)>>6))
	for _, s := range text {
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


func randFontFile() string {
  dirFIs, err := embeddedFonts.ReadDir("fonts")
  if err != nil {
    panic(err)
  }
  fonts := make([]string, 0)
  for _, dirFI := range dirFIs {
    f := filepath.Join("fonts", dirFI.Name())
    fonts = append(fonts, f)
  }

  var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
  return fonts[seededRand.Intn(len(fonts))]
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
