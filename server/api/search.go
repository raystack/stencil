package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/search"
)

func (api *API) Search(c *gin.Context) {
	ctx := c.Request.Context()
	params := c.Request.URL.Query()

	field := params.Get("field")
	namespace := params.Get("namespace")
	if field == ""{
		c.JSON(http.StatusBadRequest, gin.H{"message": "failure", "error": "field is required for search"})
	}
	
	results, err := api.SearchService.Search(ctx, &search.SearchRequest{
		Namespace: namespace,
		Field: field,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failure", "error": fmt.Sprintf("unable to search:%q", err)})
	}
	c.JSON(http.StatusOK, gin.H{"message":"success", "data": results})
}	