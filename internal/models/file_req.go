package models

type DownloadReq struct {
	AuthToken Token
	// list of tokens that each represents a file
	ObjectTokens []Token
}

// Another its name is file type
type FileExtension string

func (f FileExtension) String() string {
	return string(f)
}

type UploadReq struct {
	// authentication token. It maybe jwt or something that is agreed upon between two parties.
	AuthToken Token
	// List of tokens, each representing a file type. The value of each key is the number
	// of files we want to upload with that extension
	ObjectTypes map[FileExtension]uint
}
