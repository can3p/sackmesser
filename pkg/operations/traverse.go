package operations

import (
	"github.com/can3p/sackmesser/pkg/traverse/types"
	"github.com/pkg/errors"
)

func traverseButOne(root types.Node, path []types.PathElement) (types.Node, types.PathElement, error) {
	if len(path) < 1 {
		return root, types.PathElement{}, errors.Errorf("cannot traverse nodes with zero length path")
	}

	if len(path) == 1 {
		return root, path[0], nil
	}

	var err error
	// we do not want to traverse the last segment
	// since all the operations usually work with the parent node
	for idx := 0; idx < len(path)-1; idx++ {
		root, err = root.Visit(path[idx])

		if err != nil {
			return nil, types.PathElement{}, err
		}
	}

	return root, path[len(path)-1], nil
}
