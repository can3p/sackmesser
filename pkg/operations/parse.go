package operations

import (
	"regexp"
	"strings"

	"github.com/can3p/sackmesser/pkg/traverse/types"
	"github.com/pkg/errors"
)

var funcRegexp = regexp.MustCompile(`(\w+)\((?:(\.\w+)+)(?:\s*,\s*([^,)]+))\)`)

type Operation func(root types.Node, path []string, args ...any) error

type OpInstance struct {
	Op   Operation
	Path []string
	Args []any
}

func (op *OpInstance) Apply(root types.Node) error {
	return op.Op(root, op.Path, op.Args...)
}

var operations = map[string]Operation{
	"set": Set,
	"del": Delete,
}

// Parse is very naive at the moment, feel free to replace
// with a proper implementation. Just look how we don't handle
// anything except very simple cases!
// We should be parsing things like
// set(.field, '123') // set a string
// set(.field, some spaced value) // set a string, why not?
// set(.field, 123) // set a number
// set(.field, { a: 1 }) // assign an object to a field
// set(.field, '{ a: 1 }') // assign an string to a field
// del(.field[0].item) // delete a field
func Parse(s string) (*OpInstance, error) {
	parsed := funcRegexp.FindStringSubmatch(s)

	if len(parsed) < 3 {
		return nil, errors.Errorf("Invalid function signature")
	}

	opName := strings.ToLower(parsed[1])

	op, opExists := operations[opName]

	if !opExists {
		return nil, errors.Errorf("Operation [%s] is not supported", opName)
	}

	path := strings.Split(strings.TrimLeft(parsed[1], "."), ".")
	args := []any{}

	for _, arg := range parsed[2:] {
		parsedArg, err := ParseArg(arg)

		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse argument [%s]", parsedArg)
		}

		args = append(args, parsedArg)
	}

	return &OpInstance{
		Op:   op,
		Path: path,
		Args: args,
	}, nil
}

// implementme
func ParseArg(s string) (any, error) {
	return s, nil
}
