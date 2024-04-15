package operations

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/can3p/sackmesser/pkg/traverse/simplejson"
	"github.com/can3p/sackmesser/pkg/traverse/types"
)

func TestSetOperation(t *testing.T) {
	jstr := `{ "abc": { "def": [ 1, 2, 3 ] } }`

	examples := []struct {
		description string
		path        []types.PathElement
		arg         any
		expected    string
		isErr       bool
	}{
		{
			description: "set existing field bool",
			path:        testPath("abc"),
			arg:         true,
			expected:    `{ "abc": true }`,
		},
		// not testing integers since json parser parses everything into floats
		// by default
		{
			description: "set existing field number",
			path:        testPath("abc"),
			arg:         1234.0,
			expected:    `{ "abc": 1234.0 }`,
		},
		{
			description: "set existing field string",
			path:        testPath("abc"),
			arg:         "test",
			expected:    `{ "abc": "test" }`,
		},
		{
			description: "set existing field null",
			path:        testPath("abc"),
			arg:         nil,
			expected:    `{ "abc": null }`,
		},
		{
			description: "set existing field json",
			path:        testPath("abc"),
			arg:         map[string]any{"one": "two"},
			expected:    `{ "abc": { "one": "two" } }`,
		},
		{
			description: "set array field",
			path:        testPath("abc", "def", 0),
			arg:         "new val",
			expected:    `{ "abc": { "def": [ "new val", 2, 3 ] } }`,
		},
		{
			description: "set new field",
			path:        testPath("new field"),
			arg:         true,
			expected:    `{ "abc": { "def": [ 1, 2, 3 ] }, "new field": true }`,
		},
		{
			description: "set array field",
			path:        testPath("abc", "def", "0"),
			arg:         true,
			expected:    `{ "abc": { "def": [ true, 2, 3 ] } }`,
		},
	}

	for idx, ex := range examples {
		node := simplejson.MustParse([]byte(jstr))

		err := Set(node, ex.path, ex.arg)

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
