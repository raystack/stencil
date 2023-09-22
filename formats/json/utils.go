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
	exploreRef(jsonSchema, locationSchemaMap)
	exploreAllOf(jsonSchema, locationSchemaMap)
	exploreOneOf(jsonSchema, locationSchemaMap)
	exploreAnyOf(jsonSchema, locationSchemaMap)
	exploreProperties(jsonSchema, locationSchemaMap)
	exploreItems(jsonSchema, locationSchemaMap)
}

func exploreItems(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	if jsonSchema.Items == nil && jsonSchema.Items2020 == nil {
		return
	}
	itemSchemas := getItems(jsonSchema)
	for _, itemSchema := range itemSchemas {
		explore(itemSchema, locationSchemaMap)
	}
}

func exploreProperties(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	if jsonSchema.Properties == nil || len(jsonSchema.Properties) == 0 {
		return
	}
	for _, schema := range jsonSchema.Properties {
		explore(schema, locationSchemaMap)
	}
}

func exploreAnyOf(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	if jsonSchema.AnyOf == nil || len(jsonSchema.AnyOf) == 0 {
		return
	}
	for _, schema := range jsonSchema.AnyOf {
		explore(schema, locationSchemaMap)
	}
}

func exploreOneOf(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	if jsonSchema.OneOf == nil || len(jsonSchema.OneOf) == 0 {
		return
	}
	for _, schema := range jsonSchema.OneOf {
		explore(schema, locationSchemaMap)
	}
}

func exploreAllOf(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	if jsonSchema.AllOf == nil || len(jsonSchema.AllOf) == 0 {
		return
	}
	for _, schema := range jsonSchema.AllOf {
		explore(schema, locationSchemaMap)
	}
}

func exploreRef(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema) {
	if jsonSchema.Ref == nil {
		return
	}
	explore(jsonSchema.Ref, locationSchemaMap)
}

func elementsMatch[K comparable](arr1, arr2 []K) error {
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
			return fmt.Errorf("%v element not found in second array", element)
		} 
	}
	return nil
}

func getKeys(properties map[string]*jsonschema.Schema) []string {
	slice := make([]string, len(properties))
	for key := range properties {
		slice = append(slice, key)
	}
	return slice
}

func getDiffernce[K comparable](arr, toBeSubtracted []K) []K {
	slice := make([]K, 1)
	for _, element := range arr {
		if !contains(toBeSubtracted, element){
			slice = append(slice, element)
		} 
	}
	return slice
}

func isSubset[K comparable](superSet, subSetCandidate []K) bool {
	for _, val := range subSetCandidate {
		if !contains(superSet, val) {
			return false
		}
	}
	return true
}

func contains[K comparable](haystack []K, needle K) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}
	return false
}

func getItems(jsonSchema *jsonschema.Schema) []*jsonschema.Schema {
	schemaArr := make([]*jsonschema.Schema, 1)
	if jsonSchema.Items == nil && jsonSchema.Items2020 == nil {
		return schemaArr
	}
	var items interface{}
	if jsonSchema.Items != nil {
		items = jsonSchema.Items
	} else if jsonSchema.Items2020 != nil {
		items = jsonSchema.Items2020
	}
	itemSchema, ok := items.(*jsonschema.Schema)
	if ok {
		schemaArr = append(schemaArr, itemSchema)
		return schemaArr
	}
	itemSchemas, ok := items.([]*jsonschema.Schema)
	if ok {
		return itemSchemas
	}
	logger.Logger.Warn(fmt.Sprintf("unable to extract items schema from provided jsonschema: %s", jsonSchema.Location))
	return schemaArr
}