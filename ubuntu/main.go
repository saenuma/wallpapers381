package main

import (
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
  "github.com/pkg/errors"
  "strconv"
  "github.com/bankole7782/wallpapers381/libw381"
  "time"
)


const (
  DPI = 72.0
  SIZE = 90.0
  SPACING = 1.1
  MSIZE = 45.0
)


func main() {

  for {
    rgba := image.NewRGBA(image.Rect(0, 0, 1366, 768))

    fontBytes, err := libw381.EmbeddedFonts.ReadFile(getNextFontAddr())
    if err != nil {
      panic(err)
    }
    fontParsed, err := freetype.ParseFont(fontBytes)
    if err != nil {
      panic(err)
    }

    fontDrawer := &font.Drawer{
      Dst: rgba,
      Src: image.Black,
      Face: truetype.NewFace(fontParsed, &truetype.Options{
        Size: SIZE,
        DPI: DPI,
        Hinting: font.HintingNone,
      }),
    }

    toPrintTxt := getNextTextAddr()
    texts := wordWrap(getOutputTxt(toPrintTxt), 1366 - 130, fontDrawer)

    hex, err := colors.ParseHEX("#3C2205")
    nCR := hex.ToRGBA()
    newColor := color.RGBA{uint8(nCR.R), uint8(nCR.G), uint8(nCR.B), 255}

    hex, err = colors.ParseHEX("#F2A550")
    nCR = hex.ToRGBA()
    newColor2 := color.RGBA{uint8(nCR.R), uint8(nCR.G), uint8(nCR.B), 255}

    fg := image.NewUniform(newColor)
    bg := image.NewUniform(newColor2)

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

    fmt.Println("Changed the wallpaper @ ", time.Now().String())

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
		dd = filepath.Join(hd, "wallpapers381_data")
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


func getNextTextAddr() string {
  dirFIs, err := libw381.EmbeddedTexts.ReadDir("texts")
  if err != nil {
    panic(err)
  }

  rootPath, _ := GetRootPath()
  if DoesPathExists(filepath.Join(rootPath, "last_text.txt")) {
    rawLastText, _ := os.ReadFile(filepath.Join(rootPath, "last_text.txt"))
    number, err := strconv.Atoi(strings.TrimSpace(string(rawLastText)))
    if err != nil {
      panic(err)
    }
    toReturnNumber := number + 1
    if toReturnNumber > len(dirFIs) {
      toReturnNumber = 1
    }
    os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte(strconv.Itoa(toReturnNumber)), 0777)
    return fmt.Sprintf("texts/%d.txt", toReturnNumber)
  } else {
    os.WriteFile(filepath.Join(rootPath, "last_text.txt"), []byte("1"), 0777)
    return "texts/1.txt"
  }
}


func getNextFontAddr() string {
  dirFIs, err := libw381.EmbeddedFonts.ReadDir("fonts")
  if err != nil {
    panic(err)
  }
  fonts := make([]string, 0)
  for _, dirFI := range dirFIs {
    fonts = append(fonts, dirFI.Name())
  }

  rootPath, _ := GetRootPath()
  if DoesPathExists(filepath.Join(rootPath, "last_font.txt")) {
    rawLastFont, _ := os.ReadFile(filepath.Join(rootPath, "last_font.txt"))
    number, err := strconv.Atoi(strings.TrimSpace(string(rawLastFont)))
    if err != nil {
      panic(err)
    }
    toReturnNumber := number + 1
    if toReturnNumber > len(dirFIs) {
      toReturnNumber = 1
    }
    os.WriteFile(filepath.Join(rootPath, "last_font.txt"), []byte(strconv.Itoa(toReturnNumber)), 0777)
    return fmt.Sprintf("fonts/%s", fonts[toReturnNumber-1])
  } else {
    os.WriteFile(filepath.Join(rootPath, "last_font.txt"), []byte("1"), 0777)
    return fmt.Sprintf("fonts/%s", fonts[0])
  }
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
  rootPath, _ := GetRootPath()
  return filepath.Join(rootPath, "wallpaper.png")
}


func getOutputTxt(txtPath string) string {
  t := strings.ReplaceAll(txtPath, ".txt", "")
  t = strings.ReplaceAll(t, "texts/", "")
  textBytes, _ := libw381.EmbeddedTexts.ReadFile(txtPath)
  return t + ".  " + string(textBytes)
}
