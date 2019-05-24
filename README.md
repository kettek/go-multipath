# go-multipath
A simple library for accessing files from multiple paths.

```
var mpath multipath.Multipath

mpath.AddPath("a", multipath.FirstPriority)
mpath.AddPath("b", multipath.LastPriority)

if file, err := mpath.Open("myFile"); err != nil {
  // ...
}
defer file.Close()
// Do stuff with file
```
