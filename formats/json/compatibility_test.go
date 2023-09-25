package json

import (
	"fmt"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type compareMock struct{
	mock.Mock
}

func(m *compareMock) SchemaCompareFunc(prev, curr *jsonschema.Schema, diff *compatibilityErr){
	m.Called(prev, curr, diff)
}

func(m *compareMock) SchemaFunc(curr *jsonschema.Schema, diff *compatibilityErr){
	m.Called(curr, diff)
}

type typeCheckMock struct {
	mock.Mock
}

func (m *typeCheckMock) objectTypeChecks(prev, curr *jsonschema.Schema, diff *compatibilityErr){
	m.Called(prev, curr, diff)
}

func (m *typeCheckMock) arrayTypeChecks(prev, curr *jsonschema.Schema, diff *compatibilityErr){
	m.Called(prev, curr, diff)
}

func (m *typeCheckMock) emptyTypeChecks(prev, curr *jsonschema.Schema, diff *compatibilityErr){
	m.Called(prev, curr, diff)
}



func Test_CompareSchema_Invokes_SchemaCheck_And_Schema_CompareCheck_Expected_Number_Of_Times(t *testing.T){
	schema := initialiseSchema(t, "./testdata/compareSchemas/currSchema.json")
	prevSchema := initialiseSchema(t, "./testdata/compareSchemas/prevSchema.json")
	currMap := exploreSchema(schema)
	prevMap := exploreSchema(prevSchema)
	m := &compareMock{}
	m.On("SchemaCompareFunc", mock.Anything, mock.Anything, mock.Anything)
	m.On("SchemaFunc", mock.Anything, mock.Anything)
	diff := compareSchemas(prevMap, currMap, backwardCompatibility, []SchemaCompareCheck{m.SchemaCompareFunc}, []SchemaCheck{m.SchemaFunc})
	m.AssertNumberOfCalls(t, "SchemaCompareFunc", len(prevMap))
	m.AssertNumberOfCalls(t, "SchemaFunc", len(currMap))
	assert.Nil(t, diff) //nil because validation checks are mocked
}

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

func Test_CheckPropertyDeleted_ReturnsEmpty_When_FieldModified(t *testing.T){
	prev := initialiseSchema(t, "./testdata/propertyDeleted/prevSchema.json")
	modified := initialiseSchema(t, "./testdata/propertyDeleted/modifiedSchema.json")
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	CheckPropertyDeleted(prev, modified, diffs)
	assert.Empty(t, diffs.diffs)
}

func Test_CheckPropertyDeleted_ReturnsDiff_When_FieldDeleted(t *testing.T){
	prev := initialiseSchema(t, "./testdata/propertyDeleted/prevSchema.json")
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	CheckPropertyDeleted(prev, nil, diffs)
	assert.Equal(t, 1, len(diffs.diffs))
}

func Test_TypeCheckExecutorCorrectness(t *testing.T){
	curr := initialiseSchema(t, "./testdata/typeChecks/typeCheckSchema.json")
	objectSchema := curr.Properties["objectType"]
	arraySchema := curr.Properties["arrayType"]
	emptyTypeSchema := initialiseSchema(t, "./testdata/collection.json")
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	tcMock := &typeCheckMock{}
	tcMock.On("objectTypeChecks", objectSchema, objectSchema, diffs)
	tcMock.On("arrayTypeChecks", arraySchema, arraySchema, diffs)
	tcMock.On("emptyTypeChecks", emptyTypeSchema, emptyTypeSchema, diffs)
	spec := TypeCheckSpec{emptyTypeChecks: []SchemaCompareCheck{tcMock.emptyTypeChecks}, 
	objectTypeChecks: []SchemaCompareCheck{tcMock.objectTypeChecks}, arrayTypeChecks: []SchemaCompareCheck{tcMock.arrayTypeChecks},}
	typeCheck := TypeCheckExecutor(spec)
	typeCheck(objectSchema, objectSchema, diffs)
	tcMock.AssertNumberOfCalls(t, "objectTypeChecks", 1)
	tcMock.AssertCalled(t, "objectTypeChecks", objectSchema, objectSchema, diffs)
	tcMock.AssertNotCalled(t, "arrayTypeChecks")
	tcMock.AssertNotCalled(t, "emptyTypeChecks")
	typeCheck(arraySchema, arraySchema, diffs)
	tcMock.AssertNumberOfCalls(t, "arrayTypeChecks", 1)
	tcMock.AssertCalled(t, "arrayTypeChecks",arraySchema, arraySchema, diffs)
	tcMock.AssertNotCalled(t, "emptyTypeChecks")
	typeCheck(emptyTypeSchema, emptyTypeSchema, diffs)
	tcMock.AssertCalled(t, "emptyTypeChecks", emptyTypeSchema, emptyTypeSchema, diffs)
}

// func initSchemaCompareFunction(prev, new *jsonschema.Schema,fn SchemaCompareCheck, fn2 SchemaCheck) (map[string]*jsonschema.Schema,map[string]*jsonschema.Schema, []diffKind, []SchemaCompareCheck, []SchemaCheck){
// 	currMap := exploreSchema(new)
// 	prevMap := exploreSchema(prev)
// 	return prevMap, currMap, backwardCompatibility, []SchemaCompareCheck{fn}, []SchemaCheck{fn2}
// }

func initialiseSchema(t *testing.T, path string) *jsonschema.Schema {
	sc, err := jsonschema.Compile(path)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("failed while compiling schema: %s", path))
	}
	return sc
} 