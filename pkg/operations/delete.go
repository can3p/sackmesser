package operations

import (
	"github.com/can3p/sackmesser/pkg/traverse/types"
)

func Delete(root types.Node, path []string, args ...any) error {
	node, lastChunk, err := traverseButOne(root, path)

	if err == types.ErrFieldMissing {
		return nil
	} else if err != nil {
		return err
	}

	return node.DeleteField(lastChunk)
}
