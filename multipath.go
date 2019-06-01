package multipath

import (
	"container/list"
	"os"
	"path"
	"path/filepath"
	"sort"
)

const (
	// FirstPriority is used to add a path to the beginning of the paths list.
	FirstPriority = 0
	// LastPriority is used to add a path to the end of the paths list.
	LastPriority = -1
)

// Multipath provides a structure for opening files from an aggregated list of paths.
type Multipath struct {
	PathList list.List
}

// AddPath adds a given path to the paths list with a target priority.
func (m *Multipath) AddPath(loc string, priority int) {
	if priority == FirstPriority {
		m.PathList.PushFront(loc)
	} else if priority == LastPriority {
		m.PathList.PushBack(loc)
	} else {
		i := 0
		for e := m.PathList.Front(); e != nil; e = e.Next() {
			if i == priority {
				m.PathList.InsertBefore(loc, e)
			}
			i++
		}
	}
}

// RemovePath removes all paths matching a given path from the paths list.
func (m *Multipath) RemovePath(loc string) bool {
	for e := m.PathList.Front(); e != nil; e = e.Next() {
		p := e.Value.(string)
		if e.Value == p {
			m.PathList.Remove(e)
		}
	}
	return false
}

// Open attempts to find and open the given file path from the paths list.
func (m *Multipath) Open(name string) (*os.File, error) {
	for e := m.PathList.Front(); e != nil; e = e.Next() {
		filepath := path.Join(e.Value.(string), name)
		if file, err := os.Open(filepath); err == nil {
			return file, err
		}
	}
	return nil, os.ErrNotExist
}

// ReadFile reads the file named by filename and returns the contents.
func (m *Multipath) ReadFile(filename string) ([]byte, error) {
	file, err := m.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, fileinfo.Size())

	if _, err := file.Read(buffer); err != nil {
		return nil, err
	}
	return buffer, nil
}

// Stat attempts to find and return the FileInfo for a given file path from the paths list.
func (m *Multipath) Stat(name string) (os.FileInfo, error) {
	for e := m.PathList.Front(); e != nil; e = e.Next() {
		filepath := path.Join(e.Value.(string), name)
		if fileinfo, err := os.Stat(filepath); err == nil {
			return fileinfo, err
		}
	}
	return nil, os.ErrNotExist
}

// Lstat attempts to find and return the FileInfo for a given file path from the paths list.
func (m *Multipath) Lstat(name string) (os.FileInfo, error) {
	for e := m.PathList.Front(); e != nil; e = e.Next() {
		filepath := path.Join(e.Value.(string), name)
		if fileinfo, err := os.Lstat(filepath); err == nil {
			return fileinfo, err
		}
	}
	return nil, os.ErrNotExist
}

type walkFile struct {
	filePath string
	info     os.FileInfo
	err      error
}

// Walk walks the multipath file tree rooted at root, calling walkFn for each file or directory in the tree, including root. See https://golang.org/pkg/path/filepath/#Walk for more information.
func (m *Multipath) Walk(root string, walkFn filepath.WalkFunc) (err error) {
	filePaths := make(map[string]walkFile, 0)
	for e := m.PathList.Front(); e != nil; e = e.Next() {
		fullpath := path.Join(e.Value.(string), root)
		filepath.Walk(fullpath, func(filePath string, info os.FileInfo, err error) error {
			var localPath = filePath[len(fullpath):]
			if _, ok := filePaths[localPath]; !ok {
				filePaths[localPath] = walkFile{
					filePath: path.Clean(localPath),
					info:     info,
					err:      err,
				}
			}
			return nil
		})
	}

	var sortedFiles []string
	for k := range filePaths {
		sortedFiles = append(sortedFiles, k)
	}
	sort.Strings(sortedFiles)

	for i := range sortedFiles {
		walkFn(filePaths[sortedFiles[i]].filePath, filePaths[sortedFiles[i]].info, filePaths[sortedFiles[i]].err)
	}
	return nil
}
