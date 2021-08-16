package models

import (
	"mime/multipart"

	"github.com/odpf/stencil/server/snapshot"
)

type FileDownloadRequest struct {
	Namespace string `uri:"namespace" binding:"required"`
	Name      string `uri:"name" binding:"required"`
	Version   string `uri:"version" binding:"required,versionWithLatest"`
	FullNames []string
}

// ToSnapshot creates snapshot
func (f *FileDownloadRequest) ToSnapshot() *snapshot.Snapshot {
	s := &snapshot.Snapshot{
		Namespace: f.Namespace,
		Name:      f.Name,
	}
	if f.Version == "latest" {
		s.Latest = true
	} else {
		s.Version = f.Version
	}
	return s
}

type DescriptorUploadRequest struct {
	Namespace string                `uri:"namespace" binding:"required"`
	Name      string                `form:"name" binding:"required"`
	Version   string                `form:"version" binding:"required,version"`
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Latest    bool                  `form:"latest"`
	SkipRules []string              `form:"skiprules"`
	DryRun    bool                  `form:"dryrun"`
}

// ToSnapshot creates sanpshot
func (d *DescriptorUploadRequest) ToSnapshot() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		Namespace: d.Namespace,
		Name:      d.Name,
		Version:   d.Version,
		Latest:    d.Latest,
	}
}

type GetMetadataRequest struct {
	Namespace string `uri:"namespace" binding:"required"`
	Name      string `uri:"name" binding:"required"`
}

type MetadataUpdateRequest struct {
	Namespace string
	Name      string `json:"name" binding:"required"`
	Version   string `json:"version" binding:"required,version"`
}

type MetadataFile struct {
	Version string `json:"version"`
}
