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
  "bytes"
)


const (
  DPI = 72.0
  SIZE = 80.0
  SPACING = 1.1
)


func main() {
  fontParsed, err := freetype.ParseFont(embeddedFont)
  if err != nil {
    panic(err)
  }

  rgba := image.NewRGBA(image.Rect(0, 0, 1366, 768))
  fontDrawer := &font.Drawer{
    Dst: rgba,
    Src: image.Black,
    Face: truetype.NewFace(fontParsed, &truetype.Options{
      Size: SIZE,
      DPI: DPI,
      Hinting: font.HintingNone,
    }),
  }

  if len(os.Args) == 2 && os.Args[1] == "t" {
    dirFIs, err := embeddedTexts.ReadDir("texts")
    if err != nil {
      panic(err)
    }
    for _, dirFI := range dirFIs {
      f := filepath.Join("texts", dirFI.Name())
      textBytes, _ := embeddedTexts.ReadFile(f)
      texts := wordWrap(string(textBytes), 1366 - 200, fontDrawer)
      if len(texts) > 5 {
        panic(fmt.Sprintf("%s is more than five lines after word wrapping. Please make shorter.", f))
      }
    }
  } else {
    textBytes, err := embeddedTexts.ReadFile(randTextFile())
    if err != nil {
      panic(err)
    }
    texts := wordWrap(string(textBytes), 1366 - 130, fontDrawer)

    hex, err := colors.ParseHEX("#3C2205")
    nCR := hex.ToRGBA()
    newColor := color.RGBA{uint8(nCR.R), uint8(nCR.G), uint8(nCR.B), 255}

    // hex, err = colors.ParseHEX("#F2A550")
    // nCR = hex.ToRGBA()
    // newColor2 := color.RGBA{uint8(nCR.R), uint8(nCR.G), uint8(nCR.B), 255}
    // // #72B9AC

    bg, _, err := image.Decode(bytes.NewReader(embeddedBackground))
    if err != nil {
      panic(err)
    }
    fg := image.NewUniform(newColor)


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
  	pt := freetype.Pt(80, 50+int(c.PointToFixed(SIZE)>>6))
  	for _, s := range texts {
  		_, err = c.DrawString(s, pt)
  		if err != nil {
  			panic(err)
  		}
  		pt.Y += c.PointToFixed(SIZE * SPACING)
  	}

    // Save that RGBA image to disk.
    outPath := getOutputPath()
  	outFile, err := os.Create(outPath)
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

    fmt.Printf("Check the wallpaper at '%s'.\n", outPath)
  }

}


func randTextFile() string {
  dirFIs, err := embeddedTexts.ReadDir("texts")
  if err != nil {
    panic(err)
  }
  texts := make([]string, 0)
  for _, dirFI := range dirFIs {
    f := filepath.Join("texts", dirFI.Name())
    texts = append(texts, f)
  }

  var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
  return texts[seededRand.Intn(len(texts))]
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


func getOutputPath() string {
  hd, err := os.UserHomeDir()
	if err != nil {
	  panic("Can't get user's home directory.")
	}
	dd := os.Getenv("SNAP_USER_COMMON")
	if strings.HasPrefix(dd, filepath.Join(hd, "snap", "go")) || dd == "" {
    os.MkdirAll(filepath.Join(hd, "wallpapers381"), 0777)
		dd = filepath.Join(hd, "wallpapers381", "wallpaper.png")
	} else {
    dd = filepath.Join(dd, "wallpaper.png")
  }
  return dd
}
