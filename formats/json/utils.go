package json

import (
	"errors"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func exploreSchema(jsonSchema *jsonschema.Schema) map[string]*jsonschema.Schema {
	baseLocation := jsonSchema.Location
	exploredSchemas := make(map[string]*jsonschema.Schema, 10)
	explore(jsonSchema, exploredSchemas, baseLocation)
	return exploredSchemas
}

func explore(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema, baseLocation string) {
	_, ok := locationSchemaMap[jsonSchema.Location]
	if ok {
		return //already explored
	}
	if !strings.HasPrefix(jsonSchema.Location, baseLocation) {
		return //remote reference
	}
	locationSchemaMap[jsonSchema.Location] = jsonSchema // marking visited
	exploreRef(jsonSchema, locationSchemaMap, baseLocation)
	exploreAllOf(jsonSchema, locationSchemaMap, baseLocation)
	exploreOneOf(jsonSchema, locationSchemaMap, baseLocation)
	exploreAnyOf(jsonSchema, locationSchemaMap, baseLocation)
	exploreProperties(jsonSchema, locationSchemaMap, baseLocation)
	exploreItems(jsonSchema, locationSchemaMap, baseLocation)
	exploreAdditionalItems(jsonSchema, locationSchemaMap, baseLocation)
}

func exploreItems(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema, baseLocation string) {
	itemSchemas := getItems(jsonSchema)
	for _, itemSchema := range itemSchemas {
		explore(itemSchema, locationSchemaMap, baseLocation)
	}
}

func exploreAdditionalItems(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema, baseLocation string) {
	itemSchema := getRestOfItemsSchema(jsonSchema)
	if itemSchema != nil {
		explore(itemSchema, locationSchemaMap, baseLocation)
	}
}

func exploreProperties(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema, baseLocation string) {
	if jsonSchema.Properties == nil || len(jsonSchema.Properties) == 0 {
		return
	}
	for _, schema := range jsonSchema.Properties {
		explore(schema, locationSchemaMap, baseLocation)
	}
}

func exploreAnyOf(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema, baseLocation string) {
	if jsonSchema.AnyOf == nil || len(jsonSchema.AnyOf) == 0 {
		return
	}
	for _, schema := range jsonSchema.AnyOf {
		explore(schema, locationSchemaMap, baseLocation)
	}
}

func exploreOneOf(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema, baseLocation string) {
	if jsonSchema.OneOf == nil || len(jsonSchema.OneOf) == 0 {
		return
	}
	for _, schema := range jsonSchema.OneOf {
		explore(schema, locationSchemaMap, baseLocation)
	}
}

func exploreAllOf(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema, baseLocation string) {
	if jsonSchema.AllOf == nil || len(jsonSchema.AllOf) == 0 {
		return
	}
	for _, schema := range jsonSchema.AllOf {
		explore(schema, locationSchemaMap, baseLocation)
	}
}

func exploreRef(jsonSchema *jsonschema.Schema, locationSchemaMap map[string]*jsonschema.Schema, baseLocation string) {
	if jsonSchema.Ref == nil {
		return
	}
	explore(jsonSchema.Ref, locationSchemaMap, baseLocation)
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
		if !contains(toBeSubtracted, element) {
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

func getItems(jsSchema *jsonschema.Schema) []*jsonschema.Schema {
	if jsSchema.Draft == jsonschema.Draft2020 {
		if jsSchema.PrefixItems != nil {
			return jsSchema.PrefixItems
		}
	} else {
		schemaArr, ok := jsSchema.Items.([]*jsonschema.Schema)
		if ok {
			return schemaArr
		}
		schema, ok := jsSchema.Items.(*jsonschema.Schema)
		if ok {
			return []*jsonschema.Schema{schema}
		}
	}
	return []*jsonschema.Schema{}
}

func getRestOfItemsSchema(jsSchema *jsonschema.Schema) *jsonschema.Schema {
	if jsSchema.Draft == jsonschema.Draft2020 {
		return jsSchema.Items2020
	} else {
		schema, ok := jsSchema.AdditionalItems.(*jsonschema.Schema)
		if !ok {
			return nil
		}
		return schema
	}
}
