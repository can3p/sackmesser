package operations

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/can3p/sackmesser/pkg/traverse/simplejson"
	"github.com/can3p/sackmesser/pkg/traverse/types"
)

func TestPushOperation(t *testing.T) {
	jstr := `{ "abc": [ 1, 2, 3 ] }`

	examples := []struct {
		description string
		path        []types.PathElement
		expected    string
		arg         any
		isErr       bool
	}{
		{
			description: "push array item",
			path:        testPath("abc"),
			arg:         true,
			expected:    `{ "abc": [ 1, 2, 3, true ]}`,
		},
		{
			description: "push into scalar",
			path:        testPath("abc", 0),
			arg:         true,
			isErr:       true,
		},
	}

	for idx, ex := range examples {
		node := simplejson.MustParse([]byte(jstr))

		err := Push(node, ex.path, ex.arg)

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
