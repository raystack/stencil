package json

import (
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"

	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
)

func TestExploreJsonSchemaRecursively(t *testing.T) {
	sc, err := jsonschema.Compile("testdata/collection.json")
	assert.Nil(t, err)
	exploredMap := exploreSchema(sc)
	assert.NotEmpty(t, exploredMap)
	assert.Equal(t, 46, len(exploredMap))
}
