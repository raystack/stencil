package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CheckEnum_ForSuccess_WhenAddition_Of_Fields(t *testing.T){
	prev := initialiseSchema(t, "./testdata/enum/prev.json").Properties["roles"]
	curr := initialiseSchema(t, "./testdata/enum/curr_addition.json").Properties["roles"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkEnum(prev, curr, diffs)
	assert.Empty(t, diffs.diffs)
}

func Test_CheckEnum_ForFailure_WhenRemoval_Of_Fields(t *testing.T){
	prev := initialiseSchema(t, "./testdata/enum/prev.json").Properties["roles"]
	curr := initialiseSchema(t, "./testdata/enum/curr_removal.json").Properties["roles"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkEnum(prev, curr, diffs)
	assert.Equal(t, 1, len(diffs.diffs))
	assert.Equal(t, enumElementDeletion, diffs.diffs[0].kind)
}

func Test_CheckEnum_ForFailure_WhenEnum_Is_Removed(t *testing.T){
	prev := initialiseSchema(t, "./testdata/enum/prev.json").Properties["roles"]
	curr := initialiseSchema(t, "./testdata/enum/non_enum.json").Properties["roles"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkEnum(prev, curr, diffs)
	assert.Equal(t, 1, len(diffs.diffs))
	assert.Equal(t, enumDeletion, diffs.diffs[0].kind)
}

func Test_CheckEnum_NoPanic_WhenBothSchemaAreNonEnum(t *testing.T){
	prev := initialiseSchema(t, "./testdata/enum/non_enum.json").Properties["roles"]
	curr := initialiseSchema(t, "./testdata/enum/non_enum.json").Properties["roles"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkEnum(prev, curr, diffs)
	assert.Empty(t, diffs.diffs)
}

func Test_CheckRef_ForSuccess_WhenRefIsSame(t *testing.T){
	prev := initialiseSchema(t, "./testdata/refChange/prev.json").Properties["roleRef"]
	curr := initialiseSchema(t, "./testdata/refChange/prev.json").Properties["roleRef"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkRef(prev, curr, diffs)
	assert.Empty(t, diffs.diffs)
}

func Test_CheckRef_ForSuccess_WhenRefIsAbsentInSchemas(t *testing.T){
	prev := initialiseSchema(t, "./testdata/refChange/prev.json").Properties["roles"]
	curr := initialiseSchema(t, "./testdata/refChange/prev.json").Properties["roles"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkRef(prev, curr, diffs)
	assert.Empty(t, diffs.diffs)
}

func Test_CheckRef_ForFailure_WhenRefIsRemoved(t *testing.T){
	prev := initialiseSchema(t, "./testdata/refChange/prev.json").Properties["roleRef"]
	curr := initialiseSchema(t, "./testdata/refChange/removed.json").Properties["roleRef"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkRef(prev, curr, diffs)
	assert.Equal(t, 1, len(diffs.diffs))
	assert.Equal(t, refChanged, diffs.diffs[0].kind)
}

func Test_CheckRef_ForFailure_WhenRefIsModified(t *testing.T){
	prev := initialiseSchema(t, "./testdata/refChange/prev.json").Properties["roleRef"]
	curr := initialiseSchema(t, "./testdata/refChange/modified.json").Properties["roleRef"]
	diffs := &compatibilityErr{notAllowed: backwardCompatibility}
	checkRef(prev, curr, diffs)
	assert.Equal(t, 1, len(diffs.diffs))
	assert.Equal(t, refChanged, diffs.diffs[0].kind)
}
