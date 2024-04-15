package operations

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/can3p/sackmesser/pkg/traverse/simplejson"
	"github.com/can3p/sackmesser/pkg/traverse/types"
)

func testPath(p ...any) []types.PathElement {
	out := []types.PathElement{}

	for _, p := range p {
		str, ok := p.(string)

		if ok {
			out = append(out, types.PathElement{ObjectField: str})
			continue
		}

		intVal := p.(int)
		out = append(out, types.PathElement{ArrayIdx: intVal})
	}

	return out
}

func TestTraverseButOne(t *testing.T) {
	jstr := `{ "abc": { "def": { "cfa": "test" } } }`

	examples := []struct {
		description  string
		path         []types.PathElement
		initial      string
		expected     string
		expectedPath string
		isErr        bool
	}{
		{
			description:  "existing field",
			path:         testPath("abc", "def"),
			initial:      jstr,
			expectedPath: "def",
			expected:     `{ "def": { "cfa": "test" } }`,
		},
		{
			description:  "existing field with array in the middle",
			path:         testPath("abc", 2, "def"),
			initial:      `{ "abc": [ true, null, { "def": { "cfa": "test" } } ] }`,
			expectedPath: "def",
			expected:     `{ "def": { "cfa": "test" } }`,
		},
		{
			description:  "path with one chunk",
			path:         testPath("abc"),
			initial:      jstr,
			expectedPath: "abc",
			expected:     jstr,
		},
		{
			description: "non existant path",
			path:        testPath("ddd", "dkjk"),
			initial:     jstr,
			isErr:       true,
		},
	}

	for idx, ex := range examples {
		node := simplejson.MustParse([]byte(ex.initial))

		node, lastChunk, err := traverseButOne(node, ex.path)

		if ex.isErr {
			assert.Error(t, err, "[Ex %d - %s]", idx+1, ex.description)
			continue
		} else {
			assert.NoError(t, err, "[Ex %d - %s]", idx+1, ex.description)
		}

		expected := simplejson.MustParse([]byte(ex.expected))

		assert.Equal(t, expected.Value(), node.Value(), "[Ex %d - %s]", idx+1, ex.description)
		assert.Equal(t, ex.expectedPath, lastChunk.String(), "[Ex %d - %s]", idx+1, ex.description)
	}

}
