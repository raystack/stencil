package json

import (
	"errors"
	"fmt"

	"github.com/raystack/stencil/pkg/logger"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func exploreSchema(jsonSchema *jsonschema.Schema) map[string]*jsonschema.Schema {
	exploredSchemas := make(map[string]*jsonschema.Schema, 10)
	explore(jsonSchema, exploredSchemas)
	return exploredSchemas
}

func explore(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	_, ok := locationSchemaMap[jsonSchema.Location]
	if ok {
		return //already explored
	}
	locationSchemaMap[jsonSchema.Location] = jsonSchema // marking visited
	checkRef(jsonSchema, locationSchemaMap)
	checkAllOf(jsonSchema, locationSchemaMap)
	checkOneOf(jsonSchema, locationSchemaMap)
	checkAnyOf(jsonSchema, locationSchemaMap)
	checkProperties(jsonSchema, locationSchemaMap)
	checkItems(jsonSchema, locationSchemaMap)
}

func checkItems(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	if jsonSchema.Items == nil && jsonSchema.Items2020 == nil {
		return
	}
	var items interface{}
	if jsonSchema.Items != nil {
		items = jsonSchema.Items
	} else if jsonSchema.Items2020 != nil {
		items = jsonSchema.Items2020
	}
	itemSchema, ok := items.(*jsonschema.Schema)
	if ok {
		explore(itemSchema, locationSchemaMap)
		return
	}
	itemSchemas, ok := items.([]*jsonschema.Schema)
	if ok {
		for _, itemSchema := range itemSchemas {
			explore(itemSchema, locationSchemaMap)
		}
	}else {
		logger.Logger.Warn(fmt.Sprintf("unable to parse itemschemas to either schema or array of schemas %s", jsonSchema.Location))
	}
}

func checkProperties(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	if jsonSchema.Properties == nil || len(jsonSchema.Properties) == 0 {
		return
	}
	for _, schema := range jsonSchema.Properties {
		explore(schema, locationSchemaMap)
	}
}

func checkAnyOf(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	if jsonSchema.AnyOf == nil || len(jsonSchema.AnyOf) == 0 {
		return
	}
	for _, schema := range jsonSchema.AnyOf {
		explore(schema, locationSchemaMap)
	}
}

func checkOneOf(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	if jsonSchema.OneOf == nil || len(jsonSchema.OneOf) == 0 {
		return
	}
	for _, schema := range jsonSchema.OneOf {
		explore(schema, locationSchemaMap)
	}
}

func checkAllOf(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	if jsonSchema.AllOf == nil || len(jsonSchema.AllOf) == 0 {
		return
	}
	for _, schema := range jsonSchema.AllOf {
		explore(schema, locationSchemaMap)
	}
}

func checkRef(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	if jsonSchema.Ref == nil {
		return
	}
	explore(jsonSchema.Ref, locationSchemaMap)
}

func elementsMatch(arr1, arr2 []string) error {
	if len(arr1) != len(arr2) {
		return errors.New("count of elements do not match")
	}
	for _, element := range arr1 {
		found := false
		for _, element2 := range arr2 {
			if element == element2 {
				found = true
				break
			}
		}
		if !found {
			return errors.New(fmt.Sprintf("%s element not found in second array", element))
		} 
	}
	return nil
}