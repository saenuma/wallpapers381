package main

import (
  "math/rand"
  "time"
  "fmt"
)

func main() {
  fmt.Println(randFontFile())
}

var fonts = []string{
  "AmaticSC-Bold.ttf",
  "AnnieUseYourTelescope-Regular.ttf",
  "Itim-Regular.ttf",
  "Rancho-Regular.ttf",
}


func randFontFile() string {
  rand.Seed(time.Now().Unix())
  return fonts[rand.Intn(len(fonts))]
}
