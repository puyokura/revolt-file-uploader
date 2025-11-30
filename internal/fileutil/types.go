package fileutil

type Part struct {
	Index int    `json:"index"`
	ID    string `json:"id"`
	URL   string `json:"url"` // Optional, mainly for reference or if we want to support direct URLs
}

type FileMetadata struct {
	OriginalName string `json:"original_name"`
	TotalSize    int64  `json:"total_size"`
	Parts        []Part `json:"parts"`
}
