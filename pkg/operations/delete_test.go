package operations

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/can3p/sackmesser/pkg/traverse/simplejson"
)

func TestDeleteOperation(t *testing.T) {
	jstr := `{ "abc": { "def": [ 1, 2, 3 ] } }`

	examples := []struct {
		description string
		path        []string
		expected    string
		isErr       bool
	}{
		{
			description: "delete existing field",
			path:        []string{"abc"},
			expected:    `{}`,
		},
		{
			description: "delete missing field is fine, it was deleted already",
			path:        []string{"nonexistant"},
			expected:    jstr,
		},
	}

	for idx, ex := range examples {
		node := simplejson.MustParse([]byte(jstr))

		err := Delete(node, ex.path)

		if ex.isErr {
			assert.Error(t, err, "[Ex %d - %s]", idx+1, ex.description)
			continue
		}

		expected := simplejson.MustParse([]byte(ex.expected))

		assert.Equal(t, expected.Value(), node.Value(), "[Ex %d - %s]", idx+1, ex.description)
	}

}
