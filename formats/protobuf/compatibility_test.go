package protobuf_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/raystack/stencil/core/schema"
	"github.com/raystack/stencil/core/schema/mocks"
	"github.com/raystack/stencil/formats/protobuf"
	"github.com/stretchr/testify/assert"
)

func getCompatibilityData(t *testing.T, name string) (schema.ParsedSchema, schema.ParsedSchema) {
	t.Helper()
	rule := strings.ToLower(name)
	currentPath := filepath.Join("./testdata", rule, "current")
	currentData := getDescriptorData(t, currentPath, true)
	prevPath := filepath.Join("./testdata", rule, "previous")
	prevData := getDescriptorData(t, prevPath, true)
	current, err := protobuf.GetParsedSchema(currentData)
	assert.NoError(t, err)
	prev, err := protobuf.GetParsedSchema(prevData)
	assert.NoError(t, err)
	return current, prev
}

func TestCompatibility(t *testing.T) {
	t.Run("should return nil if no incomptible changes", func(t *testing.T) {
		current, prev := getCompatibilityData(t, "compatible")
		err := current.IsBackwardCompatible(prev)
		assert.Nil(t, err)
	})
	t.Run("backwardCompatibility", func(t *testing.T) {
		current, prev := getCompatibilityData(t, "backward")
		err := current.IsBackwardCompatible(prev)
		errMsgs := strings.Split(err.Error(), ";")
		assert.ElementsMatch(t, []string{
			`1.proto: "a.WillBeDeleted" is removed`,
			`1.proto: JSON field name changed from "nameChange" to "nameChanged"`,
			`1.proto: field "number_change" is deleted`,
			`1.proto: JSON field name changed from "numExchangeA" to "numExchangeB"`,
			`1.proto: field "num_exchange_a" kind changed from "int64" to "string"`,
			`1.proto: JSON field name changed from "numExchangeB" to "numExchangeA"`,
			`1.proto: field "num_exchange_b" kind changed from "string" to "int64"`,
			`1.proto: field "kind_change" kind changed from "string" to "int64"`,
			`1.proto: field "type_name_change" type changed from "a.TestEnum" to "a.NewTestEnum"`,
			`1.proto: field "type_message_change" type changed from "a.One" to "a.NewMessage"`,
			`1.proto: enum "a.EnumWillBeDeleted" deleted`,
			`1.proto: enum value "NAME_CHANGE" with number "3" is deleted from "a.BreakingMessage.BreackingEnum"`,
			`1.proto: enum value number for "a.BreakingMessage.NUMBER_CHANGE" changed from "1" to "2"`,
			`1.proto: previous reserved number (8) is not inclusive of current range`,
			`1.proto: previous reserved range (11, 15) is not inclusive of current range`,
			`1.proto: previous reserved name "never_existed" is removed`,
			`1.proto: field "cardinality_field" label changed from "repeated" to "optional"`,
			`2.proto: syntax changed from "proto2" to "proto3"`,
		}, errMsgs)
	})
	t.Run("backwardCompatibility return error if format does not match", func(t *testing.T) {
		current, _ := getCompatibilityData(t, "backward")
		otherSchema := &mocks.ParsedSchema{}
		otherSchema.On("Format").Return("avro")
		err := current.IsBackwardCompatible(otherSchema)
		assert.Error(t, err)
	})
	t.Run("forward", func(t *testing.T) {
		current, prev := getCompatibilityData(t, "forward")
		err := current.IsForwardCompatible(prev)
		errMsgs := strings.Split(err.Error(), ";")
		assert.ElementsMatch(t, []string{
			`1.proto: "a.WillBeDeleted" is removed`,
			`1.proto: JSON field name changed from "nameChange" to "nameChanged"`,
			`1.proto: JSON field name changed from "numExchangeA" to "numExchangeB"`,
			`1.proto: field "num_exchange_a" kind changed from "int64" to "string"`,
			`1.proto: JSON field name changed from "numExchangeB" to "numExchangeA"`,
			`1.proto: field "num_exchange_b" kind changed from "string" to "int64"`,
			`1.proto: field "kind_change" kind changed from "string" to "int64"`,
			`1.proto: field "type_name_change" type changed from "a.TestEnum" to "a.NewTestEnum"`,
			`1.proto: field "type_message_change" type changed from "a.One" to "a.NewMessage"`,
			`1.proto: enum "a.EnumWillBeDeleted" deleted`,
			`1.proto: enum value number for "a.BreakingMessage.NUMBER_CHANGE" changed from "1" to "2"`,
			`1.proto: enum value "NAME_CHANGE" with number "3" is not marked as reserved after the delete`,
			`1.proto: enum value "NAME_CHANGE" name not marked as reserved after the delete`,
			`1.proto: field "number_change" not marked as reserved after the delete`,
			`1.proto: field "delete_without_reseve" with number "8" is not marked as reserved after the delete`,
			`1.proto: field "delete_without_reseve" not marked as reserved after the delete`,
			`1.proto: field "delete_without_reseve_num" with number "9" is not marked as reserved after the delete`,
			`1.proto: field "delete_without_reseve_name" not marked as reserved after the delete`,
			`1.proto: previous reserved range (1, 6) is not inclusive of current range`,
			`1.proto: previous reserved range (8, 9) is not inclusive of current range`,
			`1.proto: previous reserved name "b" is removed`,
			`1.proto: enum value "DELETE_ENUM_WITHOUT_RESERVE_NAME" name not marked as reserved after the delete`,
			`1.proto: enum value "DELETE_ENUM_WITHOUT_RESERVE" name not marked as reserved after the delete`,
			`1.proto: enum value "DELETE_ENUM_WITHOUT_RESERVE_NUM" with number "5" is not marked as reserved after the delete`,
			`1.proto: enum value "DELETE_ENUM_WITHOUT_RESERVE" with number "4" is not marked as reserved after the delete`,
		}, errMsgs)
	})
	t.Run("forwardCompatibility return error if format does not match", func(t *testing.T) {
		current, _ := getCompatibilityData(t, "forward")
		otherSchema := &mocks.ParsedSchema{}
		otherSchema.On("Format").Return("avro")
		err := current.IsForwardCompatible(otherSchema)
		assert.Error(t, err)
	})
	t.Run("fullCompatibility", func(t *testing.T) {
		//reusing backward data
		current, prev := getCompatibilityData(t, "backward")
		err := current.IsFullCompatible(prev)
		errMsgs := strings.Split(err.Error(), ";")
		assert.ElementsMatch(t, []string{
			`1.proto: "a.WillBeDeleted" is removed`,
			`1.proto: JSON field name changed from "nameChange" to "nameChanged"`,
			`1.proto: field "number_change" is deleted`,
			`1.proto: JSON field name changed from "numExchangeA" to "numExchangeB"`,
			`1.proto: field "num_exchange_a" kind changed from "int64" to "string"`,
			`1.proto: JSON field name changed from "numExchangeB" to "numExchangeA"`,
			`1.proto: field "num_exchange_b" kind changed from "string" to "int64"`,
			`1.proto: field "kind_change" kind changed from "string" to "int64"`,
			`1.proto: field "type_name_change" type changed from "a.TestEnum" to "a.NewTestEnum"`,
			`1.proto: field "type_message_change" type changed from "a.One" to "a.NewMessage"`,
			`1.proto: enum "a.EnumWillBeDeleted" deleted`,
			`1.proto: enum value "NAME_CHANGE" with number "3" is deleted from "a.BreakingMessage.BreackingEnum"`,
			`1.proto: enum value number for "a.BreakingMessage.NUMBER_CHANGE" changed from "1" to "2"`,
			`1.proto: previous reserved number (8) is not inclusive of current range`,
			`1.proto: previous reserved range (11, 15) is not inclusive of current range`,
			`1.proto: previous reserved name "never_existed" is removed`,
			`1.proto: field "cardinality_field" label changed from "repeated" to "optional"`,
			`2.proto: syntax changed from "proto2" to "proto3"`,
		}, errMsgs)
	})
	t.Run("fullCompatibility return error if format does not match", func(t *testing.T) {
		current, _ := getCompatibilityData(t, "backward")
		otherSchema := &mocks.ParsedSchema{}
		otherSchema.On("Format").Return("avro")
		err := current.IsFullCompatible(otherSchema)
		assert.Error(t, err)
	})
}
