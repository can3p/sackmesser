package operations

import (
	"github.com/can3p/sackmesser/pkg/traverse/types"
	"github.com/pkg/errors"
)

func Push(root types.Node, path []types.PathElement, args ...any) error {
	if len(args) != 1 {
		return errors.Errorf("set operation expects one argument")
	}

	value := args[0]
	node, lastChunk, err := traverseButOne(root, path)

	if err != nil {
		return err
	}

	val, err := node.GetField(lastChunk)

	if err != nil {
		return err
	}

	typed, ok := val.([]any)
	if !ok {
		return types.ErrWrongVisit
	}

	typed = append(typed, value)

	return node.SetField(lastChunk, typed)
}
