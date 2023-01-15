# go-multipath
[![Go Reference](https://pkg.go.dev/badge/github.com/kettek/go-multipath.svg)](https://pkg.go.dev/github.com/kettek/go-multipath/v2)

A simple library for accessing files from multiple fs.FS sources as a single, unified FS. It conforms to the `FS`, `ReadDirFS`, `ReadFileFS`, `StatFS`, and `GlobFS` interfaces.

```golang
import (
  "os"

  "github.com/kettek/go-multipath/v2"
)

func main() {
  var files multipath.FS

  files.AddFS(os.DirFS("dir_a")) // has "myFile"
  files.AddFS(os.DirFS("dir_b")) // has "myFile"

  file, err := files.Open("myFile")
  if err != nil {
    panic(err)
  }
  defer file.Close()

  // Do stuff with file, which would be sourced from "dir_b".
}
```
