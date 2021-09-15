package models

import (
	"mime/multipart"
)

type FileDownloadRequest struct {
	Namespace string `uri:"namespace" binding:"required"`
	Name      string `uri:"name" binding:"required"`
	Version   string `uri:"version" binding:"required,version|eq=latest"`
	FullNames []string
}

// ToSnapshot creates snapshot
func (f *FileDownloadRequest) ToSnapshot() *Snapshot {
	s := &Snapshot{
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
func (d *DescriptorUploadRequest) ToSnapshot() *Snapshot {
	return &Snapshot{
		Namespace: d.Namespace,
		Name:      d.Name,
		Version:   d.Version,
		Latest:    d.Latest,
	}
}
