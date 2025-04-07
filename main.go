package main

import (
	"log"
	"time"
)

// type EXIF struct {}
// type IPTC struct {}

// Metadata represents the metadata of an image file, with EXIF and IPTC
// metadata as the value.
type Metadata struct {
	EXIF [][]string
	IPTC [][]string
}

// Folders represents a list of folders, with the folder names as the value.
type Folders struct {
	Folders []string
}

// Image represents an image file, with a description, a file name, a path, and
// metadata.
type Image struct {
	Description string
	File        string
	Path        string
	Metadata    Metadata
	Index       int
}

// FIXME: Rename to "Path"?
type Directory struct {
	Path string
	Name string
}

// Gallery represents a gallery, with a name, a copyright notice, a list of
// folders, a list of directories, a list of images, and a year.
type Gallery struct {
	Name        string
	Copyright   string
	Folders     []string
	Directories []Directory
	Images      []Image
	Year        int
	GalleryPath string
}

var year = time.Now().Year()

func main() {
	err := LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	err = process()
	if err != nil {
		log.Fatalf("Failed to process: %v", err)
	}
}
