package service

import (
	"context"
	"path"

	"github.com/odpf/stencil/server/models"
	"github.com/odpf/stencil/server/store"
)

//DescriptorService Interacts with backend store
type DescriptorService struct {
	Store *store.Store
}

//ListNames returns list of directories
func (d *DescriptorService) ListNames(prefixes ...string) []string {
	prefix := path.Join(prefixes...)
	paths, _ := d.Store.ListDir(prefix + "/")
	return paths
}

//ListVersions returns list of versions for specified org and name
func (d *DescriptorService) ListVersions(prefixes ...string) []string {
	prefix := path.Join(prefixes...)
	paths, _ := d.Store.ListFiles(prefix + "/")
	return paths
}

//Upload uploads the file
func (d *DescriptorService) Upload(ctx context.Context, payload *models.DescriptorPayload) error {
	filename := path.Join(payload.OrgID, payload.Name, payload.Version)
	fileReader, err := payload.File.Open()
	if err != nil {
		return err
	}
	return d.Store.Put(ctx, filename, fileReader)
}

//Download downloads the file
func (d *DescriptorService) Download(ctx context.Context, payload *models.FileMetadata) (*models.FileData, error) {
	filename := path.Join(payload.OrgID, payload.Name, payload.Version)
	data, err := d.Store.Get(ctx, filename)
	if err != nil {
		return nil, err
	}
	return &models.FileData{
		ContentLength: data.Size(),
		Reader:        data,
	}, nil
}
