package operations

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/can3p/sackmesser/pkg/traverse/simplejson"
	"github.com/can3p/sackmesser/pkg/traverse/types"
)

func TestPopOperation(t *testing.T) {
	jstr := `{ "abc": [ 1, 2, 3 ] }`

	examples := []struct {
		description string
		path        []types.PathElement
		expected    string
		isErr       bool
	}{
		{
			description: "pop array item",
			path:        testPath("abc"),
			expected:    `{ "abc": [ 1, 2 ]}`,
		},
		{
			description: "pop scalar",
			path:        testPath("abc", 0),
			isErr:       true,
		},
	}

	for idx, ex := range examples {
		node := simplejson.MustParse([]byte(jstr))

		err := Pop(node, ex.path)

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
