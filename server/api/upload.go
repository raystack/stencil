package api

import (
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/models"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// UploadDescriptor grpc handler to upload schema data with metadata information
func (a *API) UploadDescriptor(ctx context.Context, req *stencilv1.UploadDescriptorRequest) (*stencilv1.UploadDescriptorResponse, error) {
	res := &stencilv1.UploadDescriptorResponse{
		Dryrun: req.Dryrun,
	}
	s := fromProtoToSnapshot(&stencilv1.Snapshot{Namespace: req.Namespace, Name: req.Name, Version: req.Version, Latest: req.Latest})
	err := validate.StructExcept(s, "ID", "Latest")
	if err != nil {
		res.Errors = err.Error()
		return res, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := a.upload(ctx, s, req.Data, toRulesList(req.Checks), req.Dryrun); err != nil {
		res.Errors = err.Error()
		return res, err
	}
	res.Success = true
	return res, nil
}

func (a *API) upload(ctx context.Context, snapshot *models.Snapshot, data []byte, skipRules []string, dryrun bool) error {
	if ok := a.Metadata.Exists(ctx, snapshot); ok {
		return status.Error(codes.AlreadyExists, "Resource already exists")
	}
	err := a.Store.Validate(ctx, snapshot, data, skipRules)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if dryrun {
		return nil
	}
	err = a.Store.Insert(ctx, snapshot, data)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
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
