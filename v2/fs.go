package multipath

import (
	"embed"
	"io/fs"
	"os"
	"path"
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
// Other priority values will be treated as indices to insert at. If the priority is
// beyond the length of filesystems, it will be added to the end of the filesystems list.
func (m *FS) InsertFS(p fs.FS, priority int) {
	if priority == FirstPriority {
		m.filesystems = append([]fs.FS{p}, m.filesystems...)
	} else if priority == LastPriority || priority >= len(m.filesystems) {
		m.filesystems = append(m.filesystems, p)
	} else {
		m.filesystems = append(m.filesystems[:priority], append([]fs.FS{p}, m.filesystems[priority:]...)...)
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
	for _, e := range m.filesystems {
		name := name
		if _, ok := e.(embed.FS); ok {
			name = m.cleanEmbed(name)
		} else {
			name = m.Clean(name)
		}
		if d, err := e.Open(name); err == nil {
			return d, nil
		}
	}
	return nil, os.ErrNotExist
}

// ReadDir reads the named directory and returns a list of directory entries sorted by filename.
func (m *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	for _, e := range m.filesystems {
		name := name
		if _, ok := e.(embed.FS); ok {
			name = m.cleanEmbed(name)
		} else {
			name = m.Clean(name)
		}
		if d, err := fs.ReadDir(e, name); err == nil {
			return d, err
		}
	}
	return nil, os.ErrNotExist
}

// ReadFile reads the named file and returns its contents.
func (m *FS) ReadFile(name string) ([]byte, error) {
	for _, e := range m.filesystems {
		name := name
		if _, ok := e.(embed.FS); ok {
			name = m.cleanEmbed(name)
		} else {
			name = m.Clean(name)
		}
		if bytes, err := fs.ReadFile(e, name); err == nil {
			return bytes, err
		}
	}
	return nil, os.ErrNotExist
}

// Stat returns a FileInfo describing the named file from the file system.
func (m *FS) Stat(name string) (fs.FileInfo, error) {
	for _, e := range m.filesystems {
		name := name
		if _, ok := e.(embed.FS); ok {
			name = m.cleanEmbed(name)
		} else {
			name = m.Clean(name)
		}
		if info, err := fs.Stat(e, name); err == nil {
			return info, err
		}
	}
	return nil, os.ErrNotExist
}

// Glob returns the names of all files matching pattern, providing an implementation of the top-level Glob function.
func (m *FS) Glob(pattern string) ([]string, error) {
	vals := make([]string, 0)
	for _, e := range m.filesystems {
		pattern := pattern
		if _, ok := e.(embed.FS); ok {
			pattern = m.cleanEmbed(pattern)
		} else {
			pattern = m.Clean(pattern)
		}
		if matches, err := fs.Glob(e, pattern); err == nil {
			for _, m := range matches {
				matches := false
				for _, v := range vals {
					if v == m {
						matches = true
					}
				}
				if !matches {
					vals = append(vals, m)
				}
			}
		}
	}

	if len(vals) > 0 {
		return vals, nil
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
	filePaths := make(map[string]walkFile)
	for _, e := range m.filesystems {
		path := path
		if _, ok := e.(embed.FS); ok {
			path = m.cleanEmbed(path)
		} else {
			path = m.Clean(path)
		}
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

	loc = filepath.Clean(filepath.FromSlash(loc))

	if !filepath.IsAbs(loc) {
		loc = filepath.Join(string(os.PathSeparator), loc)
		loc, _ = filepath.Rel(string(os.PathSeparator), loc)
	} else {
		// Strip leading slash.
		loc = loc[1:]
	}
	return filepath.ToSlash(filepath.Clean(loc))
}

// cleanEmbed is a version of embed used only for embedded files.
func (m *FS) cleanEmbed(loc string) string {
	if loc == "" {
		return loc
	}

	loc = path.Clean("/" + loc)

	if path.IsAbs(loc) {
		// Strip leading slash.
		loc = loc[1:]
	}
	return path.Clean(loc)
}
