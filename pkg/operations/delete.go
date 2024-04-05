package operations

import (
	"github.com/can3p/sackmesser/pkg/traverse/types"
)

func Delete(root types.Node, path []string) error {
	if len(path) == 1 {
		return root.DeleteField(path[0])
	}

	node, err := root.Visit(path[0])

	if err != nil {
		return err
	}

	return Delete(node, path[1:])
}
