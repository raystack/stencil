package api

import (
    "context"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/odpf/stencil/models"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "net/http"
    "net/url"
)

func (a *API) HTTPMerge(c *gin.Context) {
    ctx := c.Request.Context()
    params := &models.FileDownloadRequest{}
    if err := c.ShouldBindUri(params); err != nil {
        c.Error(err).SetMeta(models.ErrMissingFormData)
        return
    }
    filePayload := &models.DescriptorMergeRequest{}
    if err := c.ShouldBind(filePayload); err != nil {
        c.Error(err).SetMeta(models.ErrMissingFormData)
        return
    }
    data, err := readDataFromMultiPartFile(filePayload.File)
    if err != nil {
        c.Error(err).SetMeta(models.ErrUploadInvalidFile)
        return
    }
    prevSnapshot := params.ToSnapshot()
    data, err = a.merge(ctx, prevSnapshot, data)
    if err != nil {
        c.Error(err)
        return
    }
    fileName := params.Name
    c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, fileName, url.PathEscape(fileName)))
    c.Data(http.StatusOK, "application/octet-stream", data)
}

func (a *API) merge(ctx context.Context, s *models.Snapshot, data []byte) ([]byte, error) {
    notfoundErr := status.Error(codes.NotFound, "not found")
    var prevData []byte
    st, err := a.Metadata.GetSnapshotByFields(ctx, s.Namespace, s.Name, s.Version, s.Latest)
    if err != nil {
        if err == models.ErrSnapshotNotFound {
            return data, notfoundErr
        }
        return data, status.Convert(err).Err()
    }
    prevData, err = a.Store.Get(ctx, st, []string{})
    if err != nil {
        return nil, status.Convert(err).Err()
    }
    data, err = a.Store.Merge(ctx, prevData, data)
    if err != nil {
        return nil, status.Convert(err).Err()
    }

    return data, nil
}