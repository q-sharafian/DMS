package models

import "time"

type Doc struct {
	ID ID `json:"id"`
	// The id of job position who created the document
	CreatedBy ID `json:"created_by"`
	// The id of event the document is for that
	EventID ID      `json:"event_id"`
	Context *string `json:"context"`
	// Contains path of multimedia files in the document. (If there's in the document)
	Paths []MediaPath `json:"media_paths"`
	// The time the document is created
	CreatedAt time.Time `json:"created_at"`
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
	MediaImage MediaType = iota
	MediaVideo
	MediaAudio
)
