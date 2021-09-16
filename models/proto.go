package models

// ProtobufDBFile structure to store for each file info in DB
type ProtobufDBFile struct {
	ID         int64
	SearchData *SearchData
	Data       []byte
}

// SearchData contains searchable field information
type SearchData struct {
	Path         string   `json:"path"`
	Messages     []string `json:"messages"`
	Dependencies []string `json:"dependencies"`
}
