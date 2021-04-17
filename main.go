package main

import (
  "math/rand"
  "time"
  "fmt"
  "path/filepath"
)

func main() {
  fmt.Println(randFontFile())
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
