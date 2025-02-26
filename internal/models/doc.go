package models

type Doc struct {
  ID ID `json:"id"`
  // The name of user who created the document
  CreatedBy   string `json:"created_by"`
  CreatedByID ID     `json:"created_by_id"`
  // The name of event the document is for that
  AtEvent   string `json:"at_event"`
  AtEventID ID     `json:"at_event_id"`
  Context   string `json:"content"`
  // Contains path of multimedia files in the document. (If there's in the document)
  Paths []MediaPath `json:"path"`
}

type MediaPath struct {
  Type MediaType `json:"media_type"`
  // Full path and file name (contains type too)
  Src string `json:"src"`
  // Just contains filename and its type
  FileName string `json:"file_name"`
}

type MediaType int8

const (
  Image MediaType = iota
  Video
  Audio
)
