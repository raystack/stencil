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
	if prevSchema.Ref != nil && currSchema.Ref != nil && prevSchema.Ref.Location != currSchema.Ref.Location { // check if prev and curr schema location are equivalent
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
	//  check index based schemas
	if prevSchema.Draft == jsonschema.Draft2020 {
		prevItems := prevSchema.PrefixItems
		currItems := currSchema.PrefixItems
		if prevItems != nil && currItems != nil {
			if len(prevItems) != len(currItems) {
				diffs.add(itemSchemaModification, currSchema.Location, "prev prefix items contains %d elements, current contains %d", len(prevItems), len(currItems))
			}
		}
		if prevItems == nil && currItems != nil {
			diffs.add(itemSchemaModification, currSchema.Location, "prev prefix items is absent, current contains %d", len(currItems))
		}
		if prevItems != nil && currItems == nil {
			diffs.add(itemSchemaModification, currSchema.Location, "prev prefix items contains %d elements, current contains absent", len(prevItems))
		}
	} else {
		prevItems := getItems(prevSchema)
		currItems := getItems(currSchema)
		if len(prevItems) != len(currItems) {
			diffs.add(itemSchemaModification, currSchema.Location, "prev items contains %d elements, current contains %d", len(prevItems), len(currItems))
		}
	}
}

func checkRestOfItemsSchema(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
	var prevItem, currItem *jsonschema.Schema
	var ok bool
	// check schema for remaining array elements
	if prevSchema.Draft == jsonschema.Draft2020 {
		prevItem = prevSchema.Items2020
		currItem = currSchema.Items2020
	} else {
		if prevSchema.AdditionalItems != nil {
			prevItem, ok = prevSchema.AdditionalItems.(*jsonschema.Schema)
			if !ok { // prev schema additional Items is boolean value
				if prevSchema.AdditionalItems != currSchema.AdditionalItems {
					// curr schema additional items is not equivalent
					diffs.add(itemSchemaModification, prevSchema.Location, "the value of additional items has changed")
				}
				return // since both cases equal and non equal have been evaluated.
			}
		}
		if currSchema.AdditionalItems != nil {
			currItem, ok = currSchema.AdditionalItems.(*jsonschema.Schema)
			if !ok { // curr schema is boolean
				if prevSchema.AdditionalItems == nil {
					diffs.add(itemSchemaAddition, prevSchema.Location, "additional items has been set, changes are not allowed to additional items")
				} else if prevSchema.AdditionalItems != currSchema.AdditionalItems {
					diffs.add(itemSchemaModification, prevSchema.Location, "additional items has been modified, changes are not allowed")
				}
				return
			}
		}
	}
	if prevItem == nil && currItem != nil {
		diffs.add(itemSchemaAddition, currItem.Location, "item schema cannot be added in schema changes")
	} else if prevItem != nil && currItem == nil {
		diffs.add(itemsSchemaDeletion, prevItem.Location, "items schema cannot be deleted in modification changes")
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
