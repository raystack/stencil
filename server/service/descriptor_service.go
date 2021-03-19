package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"path"
	"time"

	"github.com/odpf/stencil/server/models"
	"github.com/odpf/stencil/server/store"
)

//DescriptorService Interacts with backend store
type DescriptorService struct {
	Store *store.Store
}

//ListNames returns list of directories
func (d *DescriptorService) ListNames(prefixes ...string) ([]string, error) {
	prefix := path.Join(prefixes...)
	return d.Store.ListDir(prefix + "/")
}

//ListVersions returns list of versions for specified org and name
func (d *DescriptorService) ListVersions(prefixes ...string) ([]string, error) {
	prefix := path.Join(prefixes...)
	return d.Store.ListFiles(prefix + "/")
}

//Upload uploads the file
func (d *DescriptorService) Upload(ctx context.Context, payload *models.DescriptorPayload) error {
	orgID, name, version := payload.OrgID, payload.Name, payload.Version
	filename := path.Join(orgID, name, version)
	fileReader, err := payload.File.Open()
	if err != nil {
		return err
	}
	err = d.Store.Put(ctx, filename, fileReader)
	if err != nil {
		return err
	}
	if payload.Latest {
		return d.StoreMetadata(ctx, &models.MetadataPayload{Version: version, Name: name, OrgID: orgID})
	}
	return nil
}

//Download downloads the file
func (d *DescriptorService) Download(ctx context.Context, payload *models.FileDownload) (*models.FileData, error) {
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

//StoreMetadata stores latest version number
func (d *DescriptorService) StoreMetadata(ctx context.Context, payload *models.MetadataPayload) error {
	prefix := path.Join(payload.OrgID, payload.Name)
	metafile := path.Join(prefix, "meta.json")
	filename := path.Join(prefix, payload.Version)
	fileExists, err := d.Store.Exists(ctx, filename)
	if !fileExists {
		return models.WrapAPIError(models.ErrNotFound, err)
	}
	updated := time.Now().UTC().Format(time.RFC3339)
	fileData := &models.MetadataFile{
		Version: payload.Version,
		Updated: updated,
	}
	data, err := json.Marshal(fileData)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(data)
	err = d.Store.Put(ctx, metafile, reader)
	if err != nil {
		return err
	}
	return d.Store.Copy(ctx, filename, path.Join(prefix, "latest"))
}

//GetMetadata gets latest version number
func (d *DescriptorService) GetMetadata(ctx context.Context, payload *models.GetMetadata) (*models.MetadataFile, error) {
	filename := path.Join(payload.OrgID, payload.Name, "meta.json")
	data, err := d.Store.Get(ctx, filename)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	file := &models.MetadataFile{}
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}
