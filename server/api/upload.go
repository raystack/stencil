package api

import (
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/api/v1/genproto"
	"github.com/odpf/stencil/server/models"
	"github.com/odpf/stencil/server/snapshot"
)

// HTTPUpload http handler to schema data with metadata information
func (a *API) HTTPUpload(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	payload := &models.DescriptorUploadRequest{
		Namespace: namespace,
	}
	if err := c.ShouldBind(payload); err != nil {
		c.Error(err).SetMeta(models.ErrMissingFormData)
		return
	}
	data, err := readDataFromMultiPartFile(payload.File)
	if err != nil {
		c.Error(err).SetMeta(models.ErrUploadInvalidFile)
		return
	}
	currentSnapshot := payload.ToSnapshot()
	err = a.upload(ctx, currentSnapshot, data, payload.SkipRules, payload.DryRun)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "dryrun": payload.DryRun})
}

// Upload grpc handler to upload schema data with metadata information
func (a *API) Upload(ctx context.Context, req *genproto.UploadRequest) (*genproto.UploadResponse, error) {
	res := &genproto.UploadResponse{
		Dryrun: req.Dryrun,
	}
	s := fromProtoToSnapshot(req.Snapshot)
	if err := a.upload(ctx, s, req.Data, req.Skiprules, req.Dryrun); err != nil {
		return res, err
	}
	res.Success = true
	return res, nil
}

func (a *API) upload(ctx context.Context, snapshot *snapshot.Snapshot, data []byte, skipRules []string, dryrun bool) error {
	if ok := a.Metadata.Exists(ctx, snapshot); ok {
		return models.ErrConflict
	}
	err := a.Store.Validate(ctx, snapshot, data, skipRules)
	if err != nil {
		return models.NewAPIError(400, "", err)
	}
	if dryrun {
		return nil
	}
	err = a.Store.Insert(ctx, snapshot, data)
	if err != nil {
		return models.ErrUploadFailed
	}
	return nil
}

func readDataFromReader(reader io.ReadCloser) ([]byte, error) {
	data, err := ioutil.ReadAll(reader)
	defer func() {
		reader.Close()
	}()
	return data, err
}

func readDataFromMultiPartFile(file *multipart.FileHeader) ([]byte, error) {
	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	return readDataFromReader(fileReader)
}
