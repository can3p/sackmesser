package operations

import (
	"github.com/can3p/sackmesser/pkg/traverse/types"
	"github.com/pkg/errors"
)

func Set(root types.Node, path []string, args ...any) error {
	if len(args) != 1 {
		return errors.Errorf("add operation expects one argument")
	}

	value := args[0]

	if len(path) == 1 {
		return root.SetField(path[0], value)
	}

	node, err := root.Visit(path[0])

	if err != nil {
		return err
	}

	return Set(node, path[1:], value)
}
