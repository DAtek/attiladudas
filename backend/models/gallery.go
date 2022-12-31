package models

import "gorm.io/datatypes"

type Gallery struct {
	Id          uint
	Title       string
	Slug        string
	Description string
	Date        *datatypes.Date
	Files       []*File
	Directory   string
	Active      bool
}

type File struct {
	Id        uint
	GalleryId uint
	Filename  string
	Rank      int
}
