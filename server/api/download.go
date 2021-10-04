package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/models"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//HTTPDownload http handler to download requested schema data
func (a *API) HTTPDownload(c *gin.Context) {
	var fullNames []string
	ctx := c.Request.Context()
	fullName := c.Param("type")
	if fullName != "" {
		fullNames = []string{fullName}
	} else {
		fullNames = c.QueryArray("fullnames")
	}
	payload := models.FileDownloadRequest{
		FullNames: fullNames,
	}
	if err := c.ShouldBindUri(&payload); err != nil {
		c.Error(err).SetMeta(models.ErrMissingFormData)
		return
	}
	s := payload.ToSnapshot()
	data, err := a.download(ctx, s, payload.FullNames, payload.IsLatest())
	if err != nil {
		c.Error(err)
		return
	}
	fileName := payload.Version
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, fileName, url.PathEscape(fileName)))
	c.Data(http.StatusOK, "application/octet-stream", data)
}

// DownloadDescriptor grpc handler to download schema data
func (a *API) DownloadDescriptor(ctx context.Context, req *stencilv1.DownloadDescriptorRequest) (*stencilv1.DownloadDescriptorResponse, error) {
	payload := toFileDownloadRequest(req)
	err := validate.Struct(payload)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s := payload.ToSnapshot()
	data, err := a.download(ctx, s, req.Fullnames, payload.IsLatest())
	return &stencilv1.DownloadDescriptorResponse{Data: data}, err
}

func (a *API) download(ctx context.Context, s *models.Snapshot, fullNames []string, isLatest *bool) ([]byte, error) {
	notfoundErr := status.Error(codes.NotFound, "not found")
	var data []byte
	st, err := a.Metadata.GetSnapshotByFields(ctx, s.Namespace, s.Name, s.Version, isLatest)
	if err != nil {
		if err == models.ErrSnapshotNotFound {
			return data, notfoundErr
		}
		return data, status.Convert(err).Err()
	}
	data, err = a.Store.Get(ctx, st, fullNames)
	if err != nil {
		return data, status.Convert(err).Err()
	}
	if len(data) == 0 {
		return data, notfoundErr
	}
	return data, nil
}
