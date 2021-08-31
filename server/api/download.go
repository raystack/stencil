package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/api/v1/genproto"
	"github.com/odpf/stencil/server/models"
	"github.com/odpf/stencil/server/snapshot"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//HTTPDownload http handler to download requested schema data
func (a *API) HTTPDownload(c *gin.Context) {
	ctx := c.Request.Context()
	payload := models.FileDownloadRequest{
		FullNames: c.QueryArray("fullnames"),
	}
	if err := c.ShouldBindUri(&payload); err != nil {
		c.Error(err).SetMeta(models.ErrMissingFormData)
		return
	}
	s := payload.ToSnapshot()
	data, err := a.download(ctx, s, payload.FullNames)
	if err != nil {
		c.Error(err)
		return
	}
	fileName := payload.Version
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, fileName, url.PathEscape(fileName)))
	c.Data(http.StatusOK, "application/octet-stream", data)
}

// Download grpc handler to download schema data
func (a *API) Download(ctx context.Context, req *genproto.DownloadRequest) (*genproto.DownloadResponse, error) {
	s := fromProtoToSnapshot(req.Snapshot)
	data, err := a.download(ctx, s, req.Fullnames)
	res := &genproto.DownloadResponse{Data: data}
	return res, err
}

func (a *API) download(ctx context.Context, s *snapshot.Snapshot, fullNames []string) ([]byte, error) {
	notfoundErr := status.Error(codes.NotFound, "not found")
	var data []byte
	st, err := a.Metadata.GetSnapshotByFields(ctx, s.Namespace, s.Name, s.Version, s.Latest)
	if err != nil {
		if err == snapshot.ErrNotFound {
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
