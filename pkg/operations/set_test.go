package operations

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/can3p/sackmesser/pkg/traverse/simplejson"
)

func TestSetOperation(t *testing.T) {
	jstr := `{ "abc": { "def": [ 1, 2, 3 ] } }`

	examples := []struct {
		description string
		path        []string
		arg         any
		expected    string
		isErr       bool
	}{
		{
			description: "set existing field bool",
			path:        []string{"abc"},
			arg:         true,
			expected:    `{ "abc": true }`,
		},
		// not testing integers since json parser parses everything into floats
		// by default
		{
			description: "set existing field number",
			path:        []string{"abc"},
			arg:         1234.0,
			expected:    `{ "abc": 1234.0 }`,
		},
		{
			description: "set existing field string",
			path:        []string{"abc"},
			arg:         "test",
			expected:    `{ "abc": "test" }`,
		},
		{
			description: "set existing field null",
			path:        []string{"abc"},
			arg:         nil,
			expected:    `{ "abc": null }`,
		},
		{
			description: "set existing field json",
			path:        []string{"abc"},
			arg:         map[string]any{"one": "two"},
			expected:    `{ "abc": { "one": "two" } }`,
		},
		{
			description: "set new field",
			path:        []string{"new field"},
			arg:         true,
			expected:    `{ "abc": { "def": [ 1, 2, 3 ] }, "new field": true }`,
		},
		{
			description: "set array field",
			path:        []string{"abc", "def", "0"},
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
		}

		expected := simplejson.MustParse([]byte(ex.expected))

		assert.Equal(t, expected.Value(), node.Value(), "[Ex %d - %s]", idx+1, ex.description)
	}

}
