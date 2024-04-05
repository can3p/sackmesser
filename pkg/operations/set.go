package operations

import (
	"fmt"

	"github.com/can3p/sackmesser/pkg/traverse/types"
)

func Set(root types.Node, path []string, value any) error {
	fmt.Println("Set", path, value)
	if len(path) == 1 {
		return root.SetField(path[0], value)
	}

	node, err := root.Visit(path[0])

	if err != nil {
		return err
	}

	return Set(node, path[1:], value)
}
