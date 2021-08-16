package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/models"
	"github.com/odpf/stencil/server/snapshot"
)

//Download downloads file
func (a *API) Download(c *gin.Context) {
	ctx := c.Request.Context()
	payload := models.FileDownloadRequest{
		FullNames: c.QueryArray("fullnames"),
	}
	if err := c.ShouldBindUri(&payload); err != nil {
		c.Error(err).SetMeta(models.ErrMissingFormData)
		return
	}
	s := payload.ToSnapshot()
	st, err := a.Metadata.GetSnapshot(ctx, s.Namespace, s.Name, s.Version, s.Latest)
	if err != nil {
		if err == snapshot.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.Error(err).SetMeta(models.ErrDownloadFailed)
		return
	}
	data, err := a.Store.Get(c.Request.Context(), st, payload.FullNames)
	if err != nil {
		c.Error(err).SetMeta(models.ErrDownloadFailed)
		return
	}
	fileName := payload.Version
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, fileName, url.PathEscape(fileName)))
	c.Data(http.StatusOK, "application/octet-stream", data)
}
