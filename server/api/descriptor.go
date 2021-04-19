package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/models"
)

// ListNames lists descriptor entries
func (a *API) ListNames(c *gin.Context) {
	namespace := c.Param("namespace")
	result, err := a.Store.ListNames(namespace)
	if err != nil {
		c.Error(err).SetMeta(models.ErrUnknown)
		return
	}
	c.JSON(http.StatusOK, result)
}

// ListVersions lists version numbers for specific name
func (a *API) ListVersions(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	result, err := a.Store.ListVersions(namespace, name)
	if err != nil {
		c.Error(err).SetMeta(models.ErrUnknown)
		return
	}
	c.JSON(http.StatusOK, result)
}

//Upload uploads file
func (a *API) Upload(c *gin.Context) {
	namespace := c.Param("namespace")
	payload := models.DescriptorPayload{
		Namespace: namespace,
	}
	if err := c.ShouldBind(&payload); err != nil {
		c.Error(err).SetMeta(models.ErrMissingFormData)
		return
	}
	if err := a.Store.Upload(c.Request.Context(), &payload); err != nil {
		c.Error(err).SetMeta(models.ErrUploadFailed)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

//Download downloads file
func (a *API) Download(c *gin.Context) {
	namespace := c.Param("namespace")
	payload := models.FileDownload{
		Namespace: namespace,
	}
	if err := c.ShouldBindUri(&payload); err != nil {
		c.Error(err).SetMeta(models.ErrMissingFormData)
		return
	}
	data, err := a.Store.Download(c.Request.Context(), &payload)
	if err != nil {
		c.Error(err).SetMeta(models.ErrDownloadFailed)
		return
	}
	defer data.Reader.Close()
	fileName := c.Param("version")
	headers := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, fileName, url.PathEscape(fileName)),
	}
	c.DataFromReader(http.StatusOK, data.ContentLength, "application/octet-stream", data.Reader, headers)
}

//GetVersion return latest version number
func (a *API) GetVersion(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	data := &models.GetMetadata{
		Namespace: namespace,
		Name:      name,
	}
	version, err := a.Store.GetMetadata(c.Request.Context(), data)
	if err != nil {
		c.Error(err).SetMeta(models.ErrGetMetadataFailed)
		return
	}
	c.JSON(http.StatusOK, version)
}

//UpdateLatestVersion return latest version number
func (a *API) UpdateLatestVersion(c *gin.Context) {
	namespace := c.Param("namespace")
	payload := &models.MetadataPayload{
		Namespace: namespace,
	}
	if err := c.ShouldBind(payload); err != nil {
		c.Error(err).SetMeta(models.ErrMissingFormData)
		return
	}
	err := a.Store.StoreMetadata(c.Request.Context(), payload)
	if err != nil {
		c.Error(err).SetMeta(models.ErrMetadataUpdateFailed)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
