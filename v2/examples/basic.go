package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kettek/go-multipath/v2"
)

func main() {
	var mfs multipath.FS

	mfs.InsertFS(os.DirFS("dir_a"), multipath.FirstPriority)
	mfs.InsertFS(os.DirFS("dir_b"), multipath.LastPriority)

	file, err := mfs.Open("A")
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	file.Close()

	fmt.Printf("file 'A' contents with dir_a and dir_b:\n%s", b)

	mfs.InsertFS(os.DirFS("dir_c"), multipath.FirstPriority)

	file, err = mfs.Open("A")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err = ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("file 'A' contents with dir_a, dir_b, and dir_c:\n%s", b)
}
