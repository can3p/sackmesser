package operations

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/can3p/sackmesser/pkg/traverse/simplejson"
	"github.com/can3p/sackmesser/pkg/traverse/types"
)

func TestMergeOperation(t *testing.T) {
	jstr := `{ "abc": { "def": { "cfa": [1, 2, 3] } } }`

	examples := []struct {
		description string
		path        []types.PathElement
		arg         any
		initial     string
		expected    string
		isErr       bool
	}{
		{
			description: "add new field to the object",
			path:        testPath("abc", "def"),
			arg:         map[string]any{"added": true},
			initial:     jstr,
			expected:    `{ "abc": { "def": { "cfa": [ 1, 2, 3 ], "added": true } } }`,
		},
		{
			description: "scalar value should produce an error",
			path:        testPath("abc", "def"),
			arg:         true,
			initial:     jstr,
			isErr:       true,
		},
		{
			description: "non existant field is just set",
			path:        testPath("abc", "new field"),
			arg:         map[string]any{"added": true},
			initial:     jstr,
			expected:    `{ "abc": { "def": { "cfa": [ 1, 2, 3 ] },  "new field": { "added": true } } }`,
		},
		{
			description: "non object target field means set",
			path:        testPath("abc", "def"),
			arg:         map[string]any{"added": true},
			initial:     `{ "abc": { "def": true } }`,
			expected:    `{ "abc": { "def": { "added": true } } }`,
		},
	}

	for idx, ex := range examples {
		node := simplejson.MustParse([]byte(ex.initial))

		err := Merge(node, ex.path, ex.arg)

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
