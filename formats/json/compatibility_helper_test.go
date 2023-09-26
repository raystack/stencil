package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CheckEnum_ForSuccess_WhenAddition_Of_Fields(t *testing.T) {
	prev := initialiseSchema(t, "./testdata/enum/prev.json").Properties["roles"]
	curr := initialiseSchema(t, "./testdata/enum/curr_addition.json").Properties["roles"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkEnum(prev, curr, diffs)
	assert.Empty(t, diffs.diffs)
}

func Test_CheckEnum_ForFailure_WhenRemoval_Of_Fields(t *testing.T) {
	prev := initialiseSchema(t, "./testdata/enum/prev.json").Properties["roles"]
	curr := initialiseSchema(t, "./testdata/enum/curr_removal.json").Properties["roles"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkEnum(prev, curr, diffs)
	assert.Equal(t, 1, len(diffs.diffs))
	assert.Equal(t, enumElementDeletion, diffs.diffs[0].kind)
}

func Test_CheckEnum_ForFailure_WhenEnum_Is_Removed(t *testing.T) {
	prev := initialiseSchema(t, "./testdata/enum/prev.json").Properties["roles"]
	curr := initialiseSchema(t, "./testdata/enum/non_enum.json").Properties["roles"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkEnum(prev, curr, diffs)
	assert.Equal(t, 1, len(diffs.diffs))
	assert.Equal(t, enumDeletion, diffs.diffs[0].kind)
}

func Test_CheckEnum_NoPanic_WhenBothSchemaAreNonEnum(t *testing.T) {
	prev := initialiseSchema(t, "./testdata/enum/non_enum.json").Properties["roles"]
	curr := initialiseSchema(t, "./testdata/enum/non_enum.json").Properties["roles"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkEnum(prev, curr, diffs)
	assert.Empty(t, diffs.diffs)
}

func Test_CheckRef_ForSuccess_WhenRefIsSame(t *testing.T) {
	prev := initialiseSchema(t, "./testdata/refChange/prev.json").Properties["roleRef"]
	curr := initialiseSchema(t, "./testdata/refChange/prev.json").Properties["roleRef"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkRef(prev, curr, diffs)
	assert.Empty(t, diffs.diffs)
}

func Test_CheckRef_ForSuccess_WhenRefIsAbsentInSchemas(t *testing.T) {
	prev := initialiseSchema(t, "./testdata/refChange/prev.json").Properties["roles"]
	curr := initialiseSchema(t, "./testdata/refChange/prev.json").Properties["roles"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkRef(prev, curr, diffs)
	assert.Empty(t, diffs.diffs)
}

func Test_CheckRef_ForFailure_WhenRefIsRemoved(t *testing.T) {
	prev := initialiseSchema(t, "./testdata/refChange/prev.json").Properties["roleRef"]
	curr := initialiseSchema(t, "./testdata/refChange/removed.json").Properties["roleRef"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkRef(prev, curr, diffs)
	assert.Equal(t, 1, len(diffs.diffs))
	assert.Equal(t, refChanged, diffs.diffs[0].kind)
}

func Test_CheckRef_ForFailure_WhenRefIsModified(t *testing.T) {
	prev := initialiseSchema(t, "./testdata/refChange/prev.json").Properties["roleRef"]
	curr := initialiseSchema(t, "./testdata/refChange/modified.json").Properties["roleRef"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkRef(prev, curr, diffs)
	assert.Equal(t, 1, len(diffs.diffs))
	assert.Equal(t, refChanged, diffs.diffs[0].kind)
}

func Test_Check_AllOf_Conditions(t *testing.T) {
	prev := initialiseSchema(t, "./testdata/allOf/prev.json").Properties["roles"]
	new := initialiseSchema(t, "./testdata/allOf/modified.json").Properties["roles"]
	diffs0 := &compatibilityErr{notAllowed: backwardCompatibility}
	// check modified
	checkAllOf(prev, new, diffs0)
	assert.Equal(t, 1, len(diffs0.diffs))
	assert.Equal(t, allOfModified, diffs0.diffs[0].kind)
	// check deleted
	deleted := initialiseSchema(t, "./testdata/allOf/deleted.json").Properties["roles"]
	diffs1 := &compatibilityErr{notAllowed: backwardCompatibility}
	checkAllOf(prev, deleted, diffs1)
	assert.Equal(t, 1, len(diffs1.diffs))
	assert.Equal(t, allOfModified, diffs1.diffs[0].kind)
	// check noChange
	noChange := initialiseSchema(t, "./testdata/allOf/noChange.json").Properties["roles"]
	diffs2 := &compatibilityErr{notAllowed: backwardCompatibility}
	checkAllOf(prev, noChange, diffs2)
	assert.Empty(t, len(diffs2.diffs))
	// check addition of all of condition
	diffs3 := &compatibilityErr{notAllowed: backwardCompatibility}
	checkAllOf(deleted, prev, diffs3)
	assert.Equal(t, 1, len(diffs3.diffs))
	assert.Equal(t, allOfModified, diffs1.diffs[0].kind)
}

func Test_Check_AnyOf_Conditions(t *testing.T) {
	prev := initialiseSchema(t, "./testdata/anyOf/prev.json").Properties["roles"]
	new := initialiseSchema(t, "./testdata/anyOf/modified.json").Properties["roles"]
	diffs0 := &compatibilityErr{notAllowed: backwardCompatibility}
	// check modified
	checkAnyOf(prev, new, diffs0)
	assert.Equal(t, 1, len(diffs0.diffs))
	assert.Equal(t, anyOfModified, diffs0.diffs[0].kind)
	// check deleted
	deleted := initialiseSchema(t, "./testdata/anyOf/deleted.json").Properties["roles"]
	diffs1 := &compatibilityErr{notAllowed: backwardCompatibility}
	checkAnyOf(prev, deleted, diffs1)
	assert.Equal(t, 1, len(diffs1.diffs))
	assert.Equal(t, anyOfModified, diffs1.diffs[0].kind)
	// check noChange
	noChange := initialiseSchema(t, "./testdata/anyOf/noChange.json").Properties["roles"]
	diffs2 := &compatibilityErr{notAllowed: backwardCompatibility}
	checkAnyOf(prev, noChange, diffs2)
	assert.Empty(t, len(diffs2.diffs))
	// check addition of all of condition
	diffs3 := &compatibilityErr{notAllowed: backwardCompatibility}
	checkAnyOf(deleted, prev, diffs3)
	assert.Equal(t, 1, len(diffs3.diffs))
	assert.Equal(t, anyOfModified, diffs1.diffs[0].kind)
}

func Test_Check_OneOf_Conditions(t *testing.T) {
	prev := initialiseSchema(t, "./testdata/oneOf/prev.json").Properties["roles"]
	new := initialiseSchema(t, "./testdata/oneOf/modified.json").Properties["roles"]
	diffs0 := &compatibilityErr{notAllowed: backwardCompatibility}
	// check modified
	checkOneOf(prev, new, diffs0)
	assert.Equal(t, 1, len(diffs0.diffs))
	assert.Equal(t, oneOfModified, diffs0.diffs[0].kind)
	// check deleted
	deleted := initialiseSchema(t, "./testdata/oneOf/deleted.json").Properties["roles"]
	diffs1 := &compatibilityErr{notAllowed: backwardCompatibility}
	checkOneOf(prev, deleted, diffs1)
	assert.Equal(t, 1, len(diffs1.diffs))
	assert.Equal(t, oneOfModified, diffs1.diffs[0].kind)
	// check noChange
	noChange := initialiseSchema(t, "./testdata/oneOf/noChange.json").Properties["roles"]
	diffs2 := &compatibilityErr{notAllowed: backwardCompatibility}
	checkOneOf(prev, noChange, diffs2)
	assert.Empty(t, len(diffs2.diffs))
	// check addition of all of condition
	diffs3 := &compatibilityErr{notAllowed: backwardCompatibility}
	checkOneOf(deleted, prev, diffs3)
	assert.Equal(t, 1, len(diffs3.diffs))
	assert.Equal(t, oneOfModified, diffs1.diffs[0].kind)
}

func Test_CheckPropertyAddition_ReturnsSuccess_WhenPropertyAdded(t *testing.T) {
	prev := initialiseSchema(t, "./testdata/propertyAddition/prev.json")
	new := initialiseSchema(t, "./testdata/propertyAddition/added.json")
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkPropertyAddition(prev, new, diffs)
	// no error diffs when backward compatibility is not allowed
	assert.Empty(t, len(diffs.diffs))
	newDiff := &compatibilityErr{notAllowed: []diffKind{propertyAddition}}
	checkPropertyAddition(prev, new, newDiff)
	// diff contains element when told to record property addition
	assert.Equal(t, 1, len(newDiff.diffs))
	assert.Equal(t, propertyAddition, newDiff.diffs[0].kind)
}

func Test_CheckRequiredProperties_ReturnFailure_WhenRequiredPropertiesAdded(t *testing.T){
	prev := initialiseSchema(t, "./testdata/requiredProperties/prev.json")
	new := initialiseSchema(t, "./testdata/requiredProperties/added.json")
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkRequiredProperties(prev, new, diffs)
	assert.Equal(t, 1, len(diffs.diffs))
	assert.Equal(t, requiredFieldChanged, diffs.diffs[0].kind)
}

func Test_CheckRequiredProperties_ReturnSuccess_WhenRequiredPropertiesUnchangedAndNewPropertyAdded(t *testing.T){
	prev := initialiseSchema(t, "./testdata/requiredProperties/prev.json")
	new := initialiseSchema(t, "./testdata/requiredProperties/modified.json")
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkRequiredProperties(prev, new, diffs)
	assert.Empty(t, diffs.diffs)
}

func Test_CheckRequiredProperties_ReturnFailure_WhenRequiredPropertiesAreRemoved(t *testing.T){
	prev := initialiseSchema(t, "./testdata/requiredProperties/prev.json")
	new := initialiseSchema(t, "./testdata/requiredProperties/removed.json")
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkRequiredProperties(prev, new, diffs)
	assert.Equal(t, 1, len(diffs.diffs))
	assert.Equal(t, requiredFieldChanged, diffs.diffs[0].kind)
}
