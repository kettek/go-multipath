package main

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/kettek/go-multipath/v2"
)

func main() {
	var mfs multipath.FS

	mfs.InsertFS(os.DirFS("dir_a"), multipath.FirstPriority)
	fmt.Println("Added dir_a")
	mfs.InsertFS(os.DirFS("dir_b"), multipath.LastPriority)
	fmt.Println("Added dir_b")
	mfs.InsertFS(os.DirFS("dir_c"), multipath.FirstPriority)
	fmt.Println("Added dir_c")

	mfs.Walk(".", func(path string, d fs.DirEntry, err error) error {
		fmt.Printf("%s %+v\n", path, d.Name())
		return nil
	})

}
