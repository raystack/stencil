package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"path"
	"time"

	"github.com/odpf/stencil/server/models"
	"github.com/odpf/stencil/server/proto"
	"github.com/odpf/stencil/server/store"
)

//DescriptorService Interacts with backend store
type DescriptorService struct {
	Store        *store.Store
	ProtoService *proto.Service
}

//ListNames returns list of directories
func (d *DescriptorService) ListNames(prefixes ...string) ([]string, error) {
	return d.ProtoService.GetNames(context.Background(), prefixes[0])
}

//ListVersions returns list of versions for specified prefixes
func (d *DescriptorService) ListVersions(prefixes ...string) ([]string, error) {
	return d.ProtoService.GetVersions(context.Background(), prefixes[0], prefixes[1])
}

//Upload uploads the file
func (d *DescriptorService) Upload(ctx context.Context, payload *models.DescriptorPayload) error {
	namespace, name, version := payload.Namespace, payload.Name, payload.Version
	snapshot := &proto.Snapshot{
		Namespace: namespace,
		Name:      name,
		Version:   version,
	}
	exists := d.ProtoService.Exists(ctx, snapshot)
	if exists {
		return models.ErrConflict
	}
	data, err := readDataFromMultiPartFile(payload.File)
	if err != nil {
		return models.WrapAPIError(models.ErrUploadInvalidFile, err)
	}
	err = d.isBackwardCompatible(ctx, payload, data)
	if err != nil {
		return err
	}
	if payload.DryRun {
		return nil
	}
	return d.ProtoService.Put(ctx, &proto.Snapshot{Namespace: namespace, Name: name, Version: version, Latest: payload.Latest}, data, payload.DryRun)
}

//Download downloads the file
func (d *DescriptorService) Download(ctx context.Context, payload *models.FileDownload) (*models.FileData, error) {
	data, err := d.ProtoService.Get(ctx, &proto.Snapshot{Namespace: payload.Namespace, Name: payload.Name, Version: payload.Version}, payload.MessageFullNames)
	if err != nil {
		return nil, err
	}
	return &models.FileData{
		Data: data,
	}, nil
}

//StoreMetadata stores latest version number
func (d *DescriptorService) StoreMetadata(ctx context.Context, payload *models.MetadataPayload) error {
	prefix := path.Join(payload.Namespace, payload.Name)
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
	filename := path.Join(payload.Namespace, payload.Name, "meta.json")
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

func (d *DescriptorService) isBackwardCompatible(ctx context.Context, payload *models.DescriptorPayload, data []byte) error {
	snapshot := &proto.Snapshot{Namespace: payload.Namespace, Name: payload.Name, Version: "latest"}
	exists := d.ProtoService.Exists(ctx, snapshot)
	if !exists {
		return nil
	}
	prevData, err := d.ProtoService.Get(ctx, snapshot, []string{})
	if err != nil {
		return err
	}
	err = proto.Compare(data, prevData, payload.SkipRules)
	if err != nil {
		return models.NewAPIError(400, err.Error(), err)
	}
	return err
}
