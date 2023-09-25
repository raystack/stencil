package json

import (
	"fmt"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
)


func Test_CheckAdditionalProperties_Fails_When_Its_Partial_OpenContentModel(t *testing.T){
	schema := initialiseSchema(t, "./testdata/additionalProperties/partialOpenContent.json")
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	schemaMap := exploreSchema(schema)
	for _, schema := range schemaMap {
		CheckAdditionalProperties(schema, diffs)
	}
	assert.Equal(t, 1, len(diffs.diffs))
	assert.Equal(t, additionalPropertiesNotTrue, diffs.diffs[0].kind)
}

func Test_CheckAdditionalProperties_Fails_When_Its_ClosedContentModel(t *testing.T){
	schema := initialiseSchema(t, "./testdata/additionalProperties/closedContent.json")
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	schemaMap := exploreSchema(schema)
	for _, schema := range schemaMap {
		CheckAdditionalProperties(schema, diffs)
	}
	assert.Equal(t, 2, len(diffs.diffs))
}

func Test_CheckAdditionalProperties_Succeeds_When_Its_OpenContentModel(t *testing.T){
	schema := initialiseSchema(t, "./testdata/additionalProperties/openContent.json")
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	schemaMap := exploreSchema(schema)
	for _, schema := range schemaMap {
		CheckAdditionalProperties(schema, diffs)
	}
	assert.Empty(t, len(diffs.diffs))
}


func initialiseSchema(t *testing.T, path string) *jsonschema.Schema {
	sc, err := jsonschema.Compile(path)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("failed while compiling schema: %s", path))
	}
	return sc
} 