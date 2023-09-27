package json

import (
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func checkEnum(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
	prevEnum := prevSchema.Enum
	currEnum := currSchema.Enum
	if prevEnum == nil && currEnum != nil {
		diffs.add(enumCreation, currSchema.Location, "enum values added to existing non enum values")
	}
	if prevEnum != nil && currEnum == nil {
		diffs.add(enumDeletion, currSchema.Location, "enum was deleted")
	}
	if prevEnum != nil && currEnum != nil {
		if !isSubset(currEnum, prevEnum) {
			diffs.add(enumElementDeletion, currSchema.Location, "enum property was deleted")
		}
	}
}

func checkRef(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
	if prevSchema.Ref != nil && currSchema.Ref != nil && prevSchema.Ref.Location != currSchema.Location { // check if prev and curr schema location are equivalent
		diffs.add(refChanged, currSchema.Location, "ref for schema has been changed")
	}
	if prevSchema.Ref != nil && currSchema.Ref == nil {
		diffs.add(refChanged, currSchema.Location, "ref for schema has been removed")
	}
	if prevSchema.Ref == nil && currSchema.Ref != nil {
		diffs.add(refChanged, currSchema.Location, "ref for schema has been added")
	}
}

func checkAnyOf(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
	prevAnyOf := prevSchema.AnyOf
	currAnyOf := currSchema.AnyOf
	if prevAnyOf != nil && currAnyOf != nil {
		if len(prevAnyOf) != len(currAnyOf) {
			diffs.add(anyOfModified, currSchema.Location, "anyOf condition cannot be modified")
			return
		}
	}
	if prevAnyOf == nil && currAnyOf != nil {
		diffs.add(anyOfModified, currSchema.Location, "anyOf condition cannot created during modification of schema")
	}
	if prevAnyOf != nil && currAnyOf == nil {
		diffs.add(anyOfModified, currSchema.Location, "anyOf condition cannot be removed during modification of schema")
	}
}

func checkOneOf(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
	prevOneOf := prevSchema.OneOf
	currOneOf := currSchema.OneOf
	if prevOneOf != nil && currOneOf != nil {
		if len(prevOneOf) != len(currOneOf) {
			diffs.add(oneOfModified, currSchema.Location, "oneOf condition cannot be modified")
			return
		}
	}
	if prevOneOf == nil && currOneOf != nil {
		diffs.add(oneOfModified, currSchema.Location, "oneOf condition cannot created during modification of schema")
	}
	if prevOneOf != nil && currOneOf == nil {
		diffs.add(oneOfModified, currSchema.Location, "oneOf condition cannot be removed during modification of schema")
	}
}

func checkAllOf(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
	prevAllOf := prevSchema.AllOf
	currAllOf := currSchema.AllOf
	if prevAllOf != nil && currAllOf != nil {
		if len(prevAllOf) != len(currAllOf) {
			diffs.add(allOfModified, currSchema.Location, "allOf condition cannot be modified")
			return
		}
	}
	if prevAllOf == nil && currAllOf != nil {
		diffs.add(allOfModified, currSchema.Location, "allOf condition cannot created during modification of schema")
	}
	if prevAllOf != nil && currAllOf == nil {
		diffs.add(allOfModified, currSchema.Location, "allOf condition cannot be removed during modification of schema")
	}
}

func checkItemSchema(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
	prevItems := getItems(prevSchema)
	currItems := getItems(currSchema)
	if len(prevItems) != len(currItems) {
		diffs.add(itemSchemaModification, currSchema.Location, "prev items contains %d elements, current contains %d", len(prevItems), len(currItems))
	}
}

func checkPropertyAddition(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
	prevProperties := getKeys(prevSchema.Properties)
	currProperties := getKeys(currSchema.Properties)
	addedKeys := getDiffernce(currProperties, prevProperties)
	if len(addedKeys) > 0 {
		diffs.add(propertyAddition, currSchema.Location, "added keys: %s", strings.Join(addedKeys, ","))
	}
}

func checkRequiredProperties(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
	prevRequiredProperties := prevSchema.Required
	currReqiredProperties := currSchema.Required
	err := elementsMatch(prevRequiredProperties, currReqiredProperties)
	if err != nil {
		diffs.add(requiredFieldChanged, currSchema.Location, err.Error())
	}
}

func executeSchemaCompareCheck(prev, curr *jsonschema.Schema, diffs *compatibilityErr, checks []SchemaCompareCheck) {
	for _, check := range checks {
		check(prev, curr, diffs)
	}
}
