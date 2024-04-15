package operations

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/can3p/sackmesser/pkg/traverse/simplejson"
	"github.com/can3p/sackmesser/pkg/traverse/types"
)

func TestDeleteOperation(t *testing.T) {
	jstr := `{ "abc": { "def": [ 1, 2, 3 ] } }`

	examples := []struct {
		description string
		path        []types.PathElement
		expected    string
		isErr       bool
	}{
		{
			description: "delete existing field",
			path:        testPath("abc"),
			expected:    `{}`,
		},
		{
			description: "delete array item using object notation",
			path:        testPath("abc", "def", "1"),
			expected:    `{ "abc": { "def": [ 1, 3 ] } }`,
		},
		{
			description: "delete array item using array notation",
			path:        testPath("abc", "def", "1"),
			expected:    `{ "abc": { "def": [ 1, 3 ] } }`,
		},
		{
			description: "delete missing field is fine, it was deleted already",
			path:        testPath("nonexistant"),
			expected:    jstr,
		},
	}

	for idx, ex := range examples {
		node := simplejson.MustParse([]byte(jstr))

		err := Delete(node, ex.path)

		if ex.isErr {
			assert.Error(t, err, "[Ex %d - %s]", idx+1, ex.description)
			continue
		} else {
			assert.NoError(t, err, "[Ex %d - %s]", idx+1, ex.description)
		}

		expected := simplejson.MustParse([]byte(ex.expected))

		assert.Equal(t, expected.Value(), node.Value(), "[Ex %d - %s]", idx+1, ex.description)
	}

}
