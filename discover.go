package main

import (
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

// File represents a file on disk, with a name and a modification time.
type File struct {
	Name    string
	ModTime time.Time
}

// SubDir represents a subdirectory on disk, with a name.
type SubDir struct {
	Name string
}

// DirContent represents the content of a directory on disk.
type Dir struct {
	Name    string
	ModTime time.Time
	Files   map[string]File
	SubDirs map[string]SubDir
}

// DirMap is a map of directories on disk, with the path as the key, and a list
// of files and subdirectories as the value.
type DirMap map[string]Dir

// outputMap is a map of the output directory, with the path as the key, and the
// modification time as the value.
type outputMap map[string]time.Time

// getOutputContent returns a map of the output directory, with the path as the
// key, and the modification time as the value.
func getOutputContent() (outputMap, error) {
	results := outputMap{}
	err := filepath.WalkDir(config.Output, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileInfo, err := d.Info()
		if err != nil {
			return err
		}
		results[path] = fileInfo.ModTime()
		return nil
	})
	return results, err
}

// getOriginalContent returns a map of the originals directory, with the path as the
// key, and a list of files and subdirectories as the value.
func getOriginalContent() (DirMap, error) {
	results := DirMap{}
	err := filepath.WalkDir(config.Originals, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		fileInfo, err := d.Info()
		if err != nil {
			return err
		}
		modTime := fileInfo.ModTime()

		if d.IsDir() {
			if _, ok := results[path]; !ok {
				results[path] = Dir{
					Name:    name,
					ModTime: modTime,
					Files:   make(map[string]File),
					SubDirs: make(map[string]SubDir),
				}
			}
			if d.Name() != config.Originals {
				results[filepath.Dir(path)].SubDirs[path] = SubDir{
					Name: name,
				}
			}
		} else {
			if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") {
				results[filepath.Dir(path)].Files[path] = File{
					Name:    name,
					ModTime: modTime,
				}
			}
		}
		return nil
	})
	return results, err
}
