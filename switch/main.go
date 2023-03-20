package main

import (
  "fmt"
  "os/exec"
  "os"
  "strings"
  "path/filepath"
  "github.com/pkg/errors"
)


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


func getOutputPath() string {
  rootPath, _ := GetRootPath()
  return filepath.Join(rootPath, "wallpaper.png")
}


func main() {
  outPath := getOutputPath()

  out, err := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://" + outPath).CombinedOutput()
  if err != nil {
    fmt.Println(err)
    fmt.Println(string(out))
  }

  fmt.Println("Successfully switched to a wallpapers381's wallpaper.")
}
