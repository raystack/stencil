package models

import (
	"io"
	"mime/multipart"
)

type FileMetadata struct {
	Name    string `uri:"name" binding:"required"`
	Version string `uri:"version" binding:"required,versionWithLatest"`
	OrgID   string
}

type FileData struct {
	ContentLength int64
	ContentType   string
	Reader        io.ReadCloser
}

type DescriptorPayload struct {
	Name    string                `form:"name" binding:"required"`
	Version string                `form:"version" binding:"required,version"`
	File    *multipart.FileHeader `form:"file" binding:"required"`
	OrgID   string
}
