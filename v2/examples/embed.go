package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/kettek/go-multipath/v2"
)

//go:embed embed_dir/*
var embedFS embed.FS

func main() {
	var mfs multipath.FS

	mfs.InsertFS(os.DirFS("dir_a"), multipath.FirstPriority)
	mfs.InsertFS(os.DirFS("dir_b"), multipath.LastPriority)
	sub, err := fs.Sub(embedFS, "embed_dir")
	if err != nil {
		panic(err)
	}
	mfs.InsertFS(sub, multipath.LastPriority)

	mfs.Walk(".", func(path string, d fs.DirEntry, err error) error {
		fmt.Printf("%s %+v\n", path, d.Name())
		return nil
	})

}
