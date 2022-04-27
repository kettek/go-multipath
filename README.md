# go-multipath
[![Go Reference](https://pkg.go.dev/badge/github.com/kettek/go-multipath.svg)](https://pkg.go.dev/github.com/kettek/go-multipath)

A simple library for accessing files from multiple paths.

```
var mpath multipath.Multipath

mpath.AddPath("a", multipath.FirstPriority)
mpath.AddPath("b", multipath.LastPriority)

file, err := mpath.Open("myFile")
if err != nil {
  panic(err)
}
defer file.Close()

// Do stuff with file
```
