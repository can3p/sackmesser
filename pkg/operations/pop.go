package operations

import (
	"github.com/can3p/sackmesser/pkg/traverse/types"
)

func Pop(root types.Node, path []types.PathElement, args ...any) error {
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

	return node.SetField(lastChunk, typed[:len(typed)-1])
}
