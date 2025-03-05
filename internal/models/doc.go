package models

import "time"

type Doc struct {
	ID ID `json:"id"`
	// The id of job position who created the document
	CreatedBy ID `json:"created_by"`
	// The id of event the document is for that
	AtEvent ID     `json:"at_event"`
	Context string `json:"content"`
	// Contains path of multimedia files in the document. (If there's in the document)
	Paths []MediaPath `json:"path"`
	// The time the document is created
	At time.Time `json:"at"`
}

type MediaPath struct {
	Type MediaType `json:"media_type"`
	// Full path and file name (contains type too)
	Src string `json:"src"`
	// Just contains filename and its type
	FileName string `json:"file_name"`
}

type MediaType uint8

const (
	Image MediaType = iota
	Video
	Audio
)
