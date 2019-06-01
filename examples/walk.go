package main

import (
	"fmt"
	"github.com/kettek/go-multipath"
	"os"
)

func main() {
	var mpath multipath.Multipath

	mpath.AddPath("dir_a", multipath.FirstPriority)
	fmt.Println("Added dir_a")
	mpath.AddPath("dir_b", multipath.LastPriority)
	fmt.Println("Added dir_b")
	mpath.AddPath("dir_c", multipath.FirstPriority)
	fmt.Println("Added dir_c")

	mpath.Walk("./", func(path string, info os.FileInfo, err error) error {
		fmt.Printf("%s\n", path)
		return nil
	})

}
