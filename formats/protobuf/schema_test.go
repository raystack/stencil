package protobuf_test

import (
	"testing"

	"github.com/odpf/stencil/formats/protobuf"
	"github.com/odpf/stencil/server/schema"
	"github.com/stretchr/testify/assert"
)

func getParsedSchema(t *testing.T) schema.ParsedSchema {
	t.Helper()
	data := getDescriptorData(t, "./testdata/valid", true)
	sc, err := protobuf.GetParsedSchema(data)
	assert.NoError(t, err)
	return sc
}

func TestParsedSchema(t *testing.T) {
	t.Run("getCanonicalValue", func(t *testing.T) {
		sc := getParsedSchema(t)
		scFile := sc.GetCanonicalValue()
		assert.ElementsMatch(t, scFile.Fields, []string{"google.protobuf.Duration.seconds",
			"google.protobuf.Duration.nanos",
			"a.Test.field1",
			"a.Test.field2"})
		assert.ElementsMatch(t, []string{"google.protobuf.Duration", "a.Test"}, scFile.Types)
	})
}
