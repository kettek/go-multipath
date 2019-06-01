package main

import (
	"fmt"
	"github.com/kettek/go-multipath"
	"io/ioutil"
)

func main() {
	var mpath multipath.Multipath

	mpath.AddPath("dir_a", multipath.FirstPriority)
	mpath.AddPath("dir_b", multipath.LastPriority)

	file, err := mpath.Open("A")
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	file.Close()

	fmt.Printf("file 'A' contents with dir_a and dir_b:\n%s", b)

	mpath.AddPath("dir_c", multipath.FirstPriority)

	file, err = mpath.Open("A")
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
