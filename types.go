package main

import (
	"time"
)

// Metadata represents the metadata of an image file, with EXIF and IPTC metadata as the value.
type Metadata struct {
	EXIF [][]string
	IPTC [][]string
}

// Folders represents a list of folders, with the folder names as the value.
type Folders struct {
	Folders []string
}

type RSSItemEnclosure struct {
	URL    string
	Type   string
	Length int64
}

type RSSItem struct {
	Title       string
	Description string
	Link        string
	PubDate     string
	GUID        string
	Enclosure   RSSItemEnclosure
}

type RSSFeed struct {
	Title         string
	Description   string
	Link          string
	Copyright     string
	AtomLink      string
	Language      string
	LastBuildDate string
	Items         []RSSItem
}

// Image represents an image file, with a description, a file name, a path, and metadata.
type Image struct {
	Description string
	File        string
	Path        string
	Metadata    Metadata
	Index       int
}

// Directory represents a directory with a path and name.
// FIXME: Rename to "Path"?
type NavigationElement struct {
	Path string
	Name string
}

// Gallery represents a gallery, with metadata and content.
type Gallery struct {
	Name        string
	Copyright   string
	Folders     []string
	Navigation  []NavigationElement
	Images      []Image
	Year        int
	GalleryPath string
}

// File represents a file on disk, with a name and a modification time.
type File struct {
	Name    string
	ModTime time.Time
}

// SubDir represents a subdirectory on disk, with a name.
type SubDir struct {
	Name string
}

// Dir represents the content of a directory on disk.
type Dir struct {
	Name        string
	Path        string
	Files       map[string]File
	SubDirs     map[string]SubDir
	NeedsUpdate bool
}

// DirMap is a map of directories on disk, with the path as the key.
type DirMap map[string]Dir

// AddDir adds a new directory to the DirMap.
func (dm DirMap) AddDir(path string, name string, needsUpdate bool) {
	if _, exists := dm[path]; !exists {
		dm[path] = Dir{
			Name:        name,
			Path:        path,
			Files:       make(map[string]File),
			SubDirs:     make(map[string]SubDir),
			NeedsUpdate: needsUpdate,
		}
	}
}
