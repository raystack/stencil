package json

import (
	"fmt"

	"github.com/raystack/stencil/pkg/logger"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

const (
	_ diffKind = iota
	schemaDeleted
	incompatibleTypes
	requiredFieldChanged
	propertyAddition
	itemSchemaModification
	itemSchemaAddition
	itemsSchemaDeletion
	subSchemaTypeModification
	enumCreation
	enumDeletion
	enumElementDeletion
	refChanged
	anyOfModified
	anyOfAdded
	anyOfDeleted
	anyOfElementAdded
	anyOfElementDeleted
	oneOfAdded
	oneOfDeleted
	oneOfElementAdded
	oneOfElementDeleted
	allOfModified
	additionalPropertiesNotTrue
)

var backwardCompatibility = []diffKind{
	schemaDeleted,
	incompatibleTypes,
	requiredFieldChanged,
	itemSchemaModification,
	itemSchemaAddition,
	itemsSchemaDeletion,
	subSchemaTypeModification,
	schemaDeleted,
	incompatibleTypes,
	requiredFieldChanged,
	itemSchemaModification,
	subSchemaTypeModification,
	enumCreation,
	enumDeletion,
	enumElementDeletion,
	refChanged,
	anyOfDeleted,
	anyOfElementDeleted,
	oneOfDeleted,
	oneOfElementDeleted,
	allOfModified,
	additionalPropertiesNotTrue,
}

type SchemaCompareCheck func(prev, curr *jsonschema.Schema, err *compatibilityErr)
type SchemaCheck func(curr *jsonschema.Schema, err *compatibilityErr)

type TypeCheckSpec struct {
	emptyTypeChecks  []SchemaCompareCheck
	objectTypeChecks []SchemaCompareCheck
	arrayTypeChecks  []SchemaCompareCheck
}

var (
	emptyTypeChecks []SchemaCompareCheck = []SchemaCompareCheck{
		checkAllOf, checkAnyOf, checkOneOf, checkEnum, checkRef,
	}
	objectTypeChecks []SchemaCompareCheck = []SchemaCompareCheck{
		checkRequiredProperties, checkPropertyAddition,
	}
	/*
		Array schemas can define subschemas for each index as well as for rest of the elements.
		Hence, divided the two evaluation into two separate functions.
	*/
	arrayTypeChecks []SchemaCompareCheck = []SchemaCompareCheck{
		checkItemSchema, checkRestOfItemsSchema,
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
	/*
		enforcing open content model, in the future we can use existing additional properties schema to validate
		new properties to ensure better adherence to schema.
	*/
	if schema.AdditionalProperties != nil {
		property, ok := schema.AdditionalProperties.(bool)
		if !ok || !property {
			diffs.add(additionalPropertiesNotTrue, schema.Location, "additionalProperties need to be not defined or true for evaluation as an open content model")
		}
	}
}

func TypeCheckExecutor(spec TypeCheckSpec) SchemaCompareCheck {
	return func(prevSchema, currSchema *jsonschema.Schema, diffs *compatibilityErr) {
		if prevSchema == nil || currSchema == nil {
			return
		}
		prevTypes := prevSchema.Types
		currTypes := currSchema.Types
		err := elementsMatch(prevTypes, currTypes) // special case of integer being allowed to changed to number is not respected due to additional code complexity
		if err != nil {
			diffs.add(subSchemaTypeModification, currSchema.Location, err.Error())
			return
		}
		if len(currTypes) == 0 {
			/*
				types are not available for references and conditional schema types
				ref/holder schema
			*/
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
