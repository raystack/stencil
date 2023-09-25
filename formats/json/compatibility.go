package json

import (
	"fmt"
	"strings"

	"github.com/raystack/stencil/pkg/logger"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

const (
	_ diffKind = iota
	schemaDeleted
	incompatibleTypes
	requiredFieldChanged
	propertyAddition
	itemSchemaAddition
	subSchemaTypeModification
	enumCreation
	enumDeletion
	enumElementDeletion
	refChanged
	anyOfModified
	oneOfModified
	allOfModified
	additionalPropertiesNotTrue
)

var backwardCompatibility = []diffKind{
	schemaDeleted,
	incompatibleTypes,
	requiredFieldChanged,
	itemSchemaAddition,
	subSchemaTypeModification,
	schemaDeleted,
	incompatibleTypes,
	requiredFieldChanged,
	itemSchemaAddition,
	subSchemaTypeModification,
	enumCreation,
	enumDeletion,
	enumElementDeletion,
	refChanged,
	anyOfModified,
	oneOfModified,
	allOfModified,
	additionalPropertiesNotTrue,
}

type SchemaCompareCheck func(*jsonschema.Schema, *jsonschema.Schema, *compatibilityErr)
type SchemaCheck func(*jsonschema.Schema, *compatibilityErr)

type TypeCheckSpec struct {
	emptyTypeChecks []SchemaCompareCheck
	objectTypeChecks []SchemaCompareCheck
	arrayTypeChecks []SchemaCompareCheck
}

var (
	emptyTypeChecks []SchemaCompareCheck = []SchemaCompareCheck{
		checkAllOf, checkAnyOf, checkOneOf, checkEnum, checkRef, checkEnum,
	}
	objectTypeChecks []SchemaCompareCheck = []SchemaCompareCheck{
		checkRequiredProperties, checkFieldAddition,
	}
	arrayTypeChecks []SchemaCompareCheck = []SchemaCompareCheck{
		checkItemSchema,
	}
)

var StandardTypeChecks TypeCheckSpec = TypeCheckSpec{emptyTypeChecks, objectTypeChecks, arrayTypeChecks}

func compareSchemas(prevSchemaMap, currentSchemaMap map[string]*jsonschema.Schema, notAllowedChanges []diffKind,
	schemaCompareFuncs []SchemaCompareCheck, schemaChecks []SchemaCheck) error {
	diffs := &compatibilityErr{notAllowed: notAllowedChanges}
	for location, prevSchema := range prevSchemaMap {
		currSchema := currentSchemaMap[location]
		executeSchemaCompareCheck(prevSchema, currSchema, diffs, schemaCompareFuncs)
	}
	for _, currSchema := range currentSchemaMap {
		for _, schemaCheck := range schemaChecks {
			schemaCheck(currSchema, diffs)
		}
	}
	if diffs.isEmpty() {
		return nil
	}
	return diffs
}

func CheckPropertyDeleted(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
	if prevSchema != nil && currSchema == nil {
		diffs.add(schemaDeleted, prevSchema.Location, `property is removed`)
	}
}

func CheckAdditionalProperties(schema *jsonschema.Schema, diffs *compatibilityErr) {
	// enforcing open content model, in the future we can use existing additional properties schema to validate
	// new properties to ensure better adherence to schema.
	if schema.AdditionalProperties != nil {
		property, ok := schema.AdditionalProperties.(bool)
		if !ok || !property {
			diffs.add(additionalPropertiesNotTrue, schema.Location, "additionalProperties need to be not defined or true for evaluation as an open content model")
		}
	}

}

func TypeCheckExecutor(spec TypeCheckSpec) SchemaCompareCheck {
	return func(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
		prevTypes := prevSchema.Types
		currTypes := currSchema.Types
		err := elementsMatch(prevTypes, currTypes)
		if err != nil {
			diffs.add(subSchemaTypeModification, currSchema.Location, err.Error())
			return
		}
		if len(currTypes) == 0 {
			// types are not available for references and conditional schema types
			// ref/holder schema
			executeSchemaCompareCheck(prevSchema, currSchema, diffs, spec.emptyTypeChecks)
			return
		}
		for _, schemaTypes := range prevTypes {
			switch schemaTypes {
			case "object":
				executeSchemaCompareCheck(prevSchema, currSchema, diffs, spec.objectTypeChecks)
			case "array":
				// check item schema is same
				executeSchemaCompareCheck(prevSchema, currSchema, diffs, spec.arrayTypeChecks)
			case "integer":
				// check for validation conflicts
			case "string":
				// check validation conflicts
			case "number":
				// check validation conflicts
			case "boolean":

			case "null":

			default:
				logger.Logger.Warn(fmt.Sprintf("Unexpected type %s", schemaTypes))
			}
		}
	}
}

func checkEnum(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
	prevEnum := prevSchema.Enum
	currEnum := prevSchema.Enum
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
		for i, prevAnyOfSchema := range prevAnyOf {
			currAnyOfSchema := currAnyOf[i]
			if prevAnyOfSchema.Location != currAnyOfSchema.Location {
				diffs.add(anyOfModified, currAnyOfSchema.Location, "anyOf schema ref cannot be changed cannot be modified")
			}
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
		for i, prevOneOfSchema := range prevOneOf {
			currOneOfSchema := currOneOf[i]
			if prevOneOfSchema.Location != currOneOfSchema.Location {
				diffs.add(oneOfModified, currOneOfSchema.Location, "oneOf schema ref cannot be changed")
			}
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
		for i, prevAllOfSchema := range prevAllOf {
			currAllOfSchema := currAllOf[i]
			if prevAllOfSchema.Location != currAllOfSchema.Location {
				diffs.add(allOfModified, currAllOfSchema.Location, "allOf schema ref cannot be changed")
			}
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
		diffs.add(itemSchemaAddition, currSchema.Location, "prev items contains %d elements, current contains %d", len(prevItems), len(currItems))
	}
}

func checkFieldAddition(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
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

func executeSchemaCompareCheck(prev, curr *jsonschema.Schema, diffs *compatibilityErr, checks []SchemaCompareCheck){
	for _, check := range checks {
		check(prev, curr, diffs)
	}
}