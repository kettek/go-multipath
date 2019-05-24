package multipath

import (
	"container/list"
	"os"
	"path"
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
