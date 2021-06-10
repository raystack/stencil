package models

import (
	"io"
	"mime/multipart"
)

type FileDownload struct {
	Name      string `uri:"name" binding:"required"`
	Version   string `uri:"version" binding:"required,versionWithLatest"`
	Namespace string
}

type FileData struct {
	ContentLength int64
	ContentType   string
	Reader        io.ReadCloser
}

type DescriptorPayload struct {
	Name      string                `form:"name" binding:"required"`
	Version   string                `form:"version" binding:"required,version"`
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Latest    bool                  `form:"latest"`
	SkipRules []string              `form:"skiprules"`
	DryRun    bool                  `form:"dryrun"`
	Namespace string
}

type GetMetadata struct {
	Namespace string
	Name      string `uri:"name"`
}

type MetadataPayload struct {
	Version   string `form:"version" json:"version" binding:"required,version"`
	Namespace string
	Name      string `json:"name" binding:"required"`
}

type MetadataFile struct {
	Version string `json:"version"`
	Updated string `json:"updated"`
}
