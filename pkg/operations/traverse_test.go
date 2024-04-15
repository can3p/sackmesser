package operations

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/can3p/sackmesser/pkg/traverse/simplejson"
)

func TestTraverseButOne(t *testing.T) {
	jstr := `{ "abc": { "def": { "cfa": "test" } } }`

	examples := []struct {
		description  string
		path         []string
		initial      string
		expected     string
		expectedPath string
		isErr        bool
	}{
		{
			description:  "existing field",
			path:         []string{"abc", "def"},
			initial:      jstr,
			expectedPath: "def",
			expected:     `{ "def": { "cfa": "test" } }`,
		},
		{
			description:  "path with one chunk",
			path:         []string{"abc"},
			initial:      jstr,
			expectedPath: "abc",
			expected:     jstr,
		},
		{
			description: "non existant path",
			path:        []string{"ddd", "dkjk"},
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
		assert.Equal(t, ex.expectedPath, lastChunk, "[Ex %d - %s]", idx+1, ex.description)
	}

}
