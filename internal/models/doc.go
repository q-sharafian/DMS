package models

import "encoding/json"

type Doc struct {
	ID ID `json:"id" example:"20354d7a-e4fe-47af-8ff6-187bca92f3f9"`
	// The id of job position who created the document
	CreatedBy ID `json:"created_by" example:"54a79030f-0685-49d1-bbdd-31ab1b4c1613" validate:"required"`
	// The id of event the document is for that
	EventID ID      `json:"event_id" example:"32a79030f-0685-49d1-bbdd-31ab1b4c1613" validate:"required"`
	Context *string `json:"context" example:"some context"`
	// Contains path of multimedia files in the document. (If there's in the document)
	Paths []MediaPath `json:"media_paths"`
	// The time the document is created. It's in UTC time zone and Unix timestamp. (in seconds)
	CreatedAt int64 `json:"created_at" example:"1641011200"`
}

type MediaPath struct {
	Type MediaType `json:"media_type"`
	// Full path and file name (contains type too)
	Src string `json:"src"`
	// Just contains filename and its type
	FileName string `json:"file_name"`
}

func (s Doc) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        string      `json:"id"`
		CreatedBy string      `json:"created_by"`
		EventID   string      `json:"event_id"`
		Context   *string     `json:"context"`
		Paths     []MediaPath `json:"media_paths"`
		CreatedAt int64       `json:"created_at"`
	}{
		ID:        s.ID.ToString(),
		CreatedBy: s.CreatedBy.ToString(),
		EventID:   s.EventID.ToString(),
		Context:   s.Context,
		Paths:     s.Paths,
		CreatedAt: s.CreatedAt,
	})
}

type MediaType uint8

const (
	MediaImage MediaType = iota
	MediaVideo
	MediaAudio
)
