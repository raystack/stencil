package api

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/models"
)

//Upload uploads file
func (a *API) Upload(c *gin.Context) {
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
	if ok := a.Metadata.Exists(ctx, currentSnapshot); ok {
		c.JSON(http.StatusConflict, gin.H{"message": "Resource already exist"})
		return
	}
	err = a.Store.Validate(ctx, currentSnapshot, data, payload.SkipRules)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if payload.DryRun {
		c.JSON(http.StatusOK, gin.H{"message": "success", "dryrun": "true"})
		return
	}
	err = a.Store.Insert(ctx, currentSnapshot, data)
	if err != nil {
		c.Error(err).SetMeta(models.ErrUploadFailed)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
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
