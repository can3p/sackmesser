package operations

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/can3p/sackmesser/pkg/traverse/types"
)

func TestParse(t *testing.T) {
	ex := []struct {
		description  string
		input        string
		isError      bool
		ExpectedOp   string
		ExpectedArgs []any
		ExpectedPath []types.PathElement
	}{
		{
			description:  "test nested path",
			input:        "set(field.another_field, true)",
			ExpectedOp:   "set",
			ExpectedPath: testPath("field", "another_field"),
			ExpectedArgs: []any{true},
		},
		{
			description:  "test nested path with array lookup",
			input:        `set(field[0].another_field[8][9]."one more", true)`,
			ExpectedOp:   "set",
			ExpectedPath: testPath("field", 0, "another_field", 8, 9, "one more"),
			ExpectedArgs: []any{true},
		},
		{
			description:  "test path with strings",
			input:        "set(field.\"another field\", true)",
			ExpectedOp:   "set",
			ExpectedPath: testPath("field", "another field"),
			ExpectedArgs: []any{true},
		},
		{
			description:  "test boolean",
			input:        "set(field, true)",
			ExpectedOp:   "set",
			ExpectedPath: testPath("field"),
			ExpectedArgs: []any{true},
		},
		{
			description:  "test int",
			input:        "set(field, 12345)",
			ExpectedOp:   "set",
			ExpectedPath: testPath("field"),
			ExpectedArgs: []any{12345},
		},
		{
			description:  "test string with single quotes",
			input:        `set(field, '123"   45')`,
			ExpectedOp:   "set",
			ExpectedPath: testPath("field"),
			ExpectedArgs: []any{"123\"   45"},
		},
		{
			description:  "test string with double quotes",
			input:        `set(field, "123'   45")`,
			ExpectedOp:   "set",
			ExpectedPath: testPath("field"),
			ExpectedArgs: []any{"123'   45"},
		},
		{
			description:  "test string with back ticks",
			input:        "set(field, `123'\"   45`)",
			ExpectedOp:   "set",
			ExpectedPath: testPath("field"),
			ExpectedArgs: []any{`123'"   45`},
		},
		{
			description:  "test string with single quotes",
			input:        `set(field, '123"   45')`,
			ExpectedOp:   "set",
			ExpectedPath: testPath("field"),
			ExpectedArgs: []any{"123\"   45"},
		},
		{
			description:  "test bare word without quotes",
			input:        `set(field, awesome)`,
			ExpectedOp:   "set",
			ExpectedPath: testPath("field"),
			ExpectedArgs: []any{"awesome"},
		},
		{
			description:  "test null",
			input:        "set(field, null)",
			ExpectedOp:   "set",
			ExpectedPath: testPath("field"),
			ExpectedArgs: []any{nil},
		},
		{
			description:  "test json",
			input:        `set(field, { "a": true })`,
			ExpectedOp:   "set",
			ExpectedPath: testPath("field"),
			ExpectedArgs: []any{map[string]any{"a": true}},
		},
		{
			description:  "test one more json",
			input:        `set(field, { "abc": [1,2, false] })`,
			ExpectedOp:   "set",
			ExpectedPath: testPath("field"),
			ExpectedArgs: []any{map[string]any{"abc": []any{float64(1), float64(2), false}}},
		},
	}

	parser := NewParser()

	for idx, ex := range ex {
		parsed, err := parser.Parse(ex.input)

		if ex.isError != (err != nil) {
			if ex.isError {
				t.Errorf("[%d - %s] expected an error, but got none", idx+1, ex.description)
			} else {
				t.Errorf("[%d - %s] expected no error, but one: %s", idx+1, ex.description, err.Error())
			}
		}

		if err != nil {
			continue
		}

		if ex.ExpectedOp != parsed.Name {
			t.Errorf("[%d - %s] expected op %s, but got %s", idx+1, ex.description, ex.ExpectedOp, parsed.Name)
			continue
		}

		assert.Equal(t, ex.ExpectedPath, parsed.Path, "[%d - %s] expected path %s, but got %s", idx+1, ex.description, ex.ExpectedPath, parsed.Path)

		assert.Equal(t, ex.ExpectedArgs, parsed.Args, "[%d - %s] arguments mismatch", idx+1, ex.description)
	}
}
