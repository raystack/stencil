package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/models"
	"github.com/odpf/stencil/server/snapshot"
)

// ListNames lists descriptor entries
func (a *API) ListNames(c *gin.Context) {
	namespace := c.Param("namespace")
	result, err := a.Metadata.ListNames(c.Request.Context(), namespace)
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
	result, err := a.Metadata.ListVersions(c.Request.Context(), namespace, name)
	if err != nil {
		c.Error(err).SetMeta(models.ErrUnknown)
		return
	}
	c.JSON(http.StatusOK, result)
}

//GetLatestVersion return latest version number
func (a *API) GetLatestVersion(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	snapshot, err := a.Metadata.GetSnapshot(c.Request.Context(), namespace, name, "", true)
	if err != nil {
		c.Error(err).SetMeta(models.ErrGetMetadataFailed)
		return
	}
	c.JSON(http.StatusOK, gin.H{"version": snapshot.Version})
}

//UpdateLatestVersion return latest version number
func (a *API) UpdateLatestVersion(c *gin.Context) {
	namespace := c.Param("namespace")
	payload := &models.MetadataUpdateRequest{
		Namespace: namespace,
	}
	if err := c.ShouldBind(payload); err != nil {
		c.Error(err).SetMeta(models.ErrMissingFormData)
		return
	}
	err := a.Metadata.UpdateLatestVersion(c.Request.Context(), &snapshot.Snapshot{
		Namespace: namespace,
		Name:      payload.Name,
		Version:   payload.Version,
	})
	if err != nil {
		c.Error(err).SetMeta(models.ErrMetadataUpdateFailed)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
