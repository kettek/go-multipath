package multipath

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// FS is a container for multiple fs.FS conforming items.
type FS struct {
	filesystems []fs.FS
}

// These are our priorities.
const (
	FirstPriority = 0
	LastPriority  = -1
)

// AddFS adds the given fs to the beginning of the filesystems list.
func (m *FS) AddFS(p fs.FS) {
	m.InsertFS(p, FirstPriority)
}

// InsertFS adds the given FS to the filesystems list with a given priority.
// If priority is FirstPriority, it is added to the beginning of the filesystems list.
// If it is LastPriority, it is added at the end of the filesystems list.
// Other priority values will be treated as indices to insert at.
func (m *FS) InsertFS(p fs.FS, priority int) {
	if priority == FirstPriority {
		m.filesystems = append([]fs.FS{p}, m.filesystems...)
	} else if priority == LastPriority {
		m.filesystems = append(m.filesystems, p)
	} else {
		for i := 0; i < len(m.filesystems); i++ {
			if i == priority {
				m.filesystems = append(m.filesystems[:i], append([]fs.FS{p}, m.filesystems[i:]...)...)
			}
		}
	}
}

// RemoveFS removes the given FS from the filesystems list.
func (m *FS) RemoveFS(p fs.FS) bool {
	for i, e := range m.filesystems {
		if e == p {
			m.filesystems = append(m.filesystems[:i], m.filesystems[i+1:]...)
			return true
		}
	}
	return false
}

// Open opens the named file.
func (m *FS) Open(name string) (fs.File, error) {
	name = m.Clean(name)
	for _, e := range m.filesystems {
		if e, ok := e.(fs.FS); ok {
			if d, err := e.Open(name); err == nil {
				return d, nil
			}
		}
	}
	return nil, os.ErrNotExist
}

// ReadDir reads the named directory and returns a list of directory entries sorted by filename.
func (m *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	name = m.Clean(name)
	for _, e := range m.filesystems {
		if e, ok := e.(fs.ReadDirFS); ok {
			if d, err := e.ReadDir(name); err == nil {
				return d, nil
			}
		}
	}
	return nil, os.ErrNotExist
}

// ReadFile reads the named file and returns its contents.
func (m *FS) ReadFile(name string) ([]byte, error) {
	name = m.Clean(name)
	for _, e := range m.filesystems {
		if e, ok := e.(fs.ReadFileFS); ok {
			if f, err := e.ReadFile(name); err == nil {
				return f, nil
			}
		}
	}
	return nil, os.ErrNotExist
}

// Stat returns a FileInfo describing the named file from the file system.
func (m *FS) Stat(name string) (fs.FileInfo, error) {
	name = m.Clean(name)
	for _, e := range m.filesystems {
		if e, ok := e.(fs.StatFS); ok {
			if s, err := e.Stat(name); err == nil {
				return s, nil
			}
		}
	}
	return nil, os.ErrNotExist
}

// Glob returns the names of all files matching pattern, providing an implementation of the top-level Glob function.
func (m *FS) Glob(pattern string) ([]string, error) {
	pattern = m.Clean(pattern) // Hmm... this might not work right.
	for _, e := range m.filesystems {
		if e, ok := e.(fs.GlobFS); ok {
			if s, err := e.Glob(pattern); err == nil {
				return s, nil
			}
		}
	}
	return nil, os.ErrNotExist
}

type walkFile struct {
	filePath string
	d        fs.DirEntry
	err      error
}

// Walk traverses the directory structure, calling wallkFn on each.
func (m *FS) Walk(path string, walkFn fs.WalkDirFunc) (err error) {
	path = m.Clean(path)
	filePaths := make(map[string]walkFile)
	for _, e := range m.filesystems {
		fs.WalkDir(e, path, func(path string, d fs.DirEntry, err error) error {
			if _, ok := filePaths[path]; !ok {
				filePaths[path] = walkFile{
					filePath: path,
					d:        d,
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
		walkFn(filePaths[sortedFiles[i]].filePath, filePaths[sortedFiles[i]].d, filePaths[sortedFiles[i]].err)
	}
	return nil
}

// Clean cleans a path to remove access to unsafe directories.
func (m *FS) Clean(loc string) string {
	if loc == "" {
		return loc
	}
	loc = filepath.Clean(loc)

	if !filepath.IsAbs(loc) {
		loc = filepath.Join(string(os.PathSeparator), loc)
		loc, _ = filepath.Rel(string(os.PathSeparator), loc)
	} else {
		// Strip leading slash.
		loc = loc[1:]
	}
	return filepath.Clean(loc)
}
