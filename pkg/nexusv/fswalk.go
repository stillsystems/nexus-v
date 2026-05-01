package nexusv

import (
	"io/fs"
)

// fsWalk walks the embedded template filesystem starting at root
// and calls fn for every file and directory encountered.
func fsWalk(root string, fn func(path string, isDir bool) error) error {
	return fs.WalkDir(templateFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == root {
			return nil
		}

		return fn(path, d.IsDir())
	})
}
