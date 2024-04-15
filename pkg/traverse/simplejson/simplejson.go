package simplejson

import (
	"encoding/json"

	"github.com/can3p/sackmesser/pkg/traverse/simpleobject"
	"github.com/can3p/sackmesser/pkg/traverse/types"
)

type jnode struct {
	types.Node
}

func (n *jnode) Serialize() ([]byte, error) {
	return json.MarshalIndent(n.Value(), "", "  ")
}

func Parse(b []byte) (types.RootNode, error) {
	var j any

	if err := json.Unmarshal(b, &j); err != nil {
		return nil, err
	}

	return &jnode{
		simpleobject.FromValue(j),
	}, nil
}

func MustParse(b []byte) types.Node {
	n, err := Parse(b)

	if err != nil {
		panic(err)
	}

	return n
}

func FromNode(n types.Node) types.RootNode {
	return &jnode{
		simpleobject.FromNode(n),
	}
}
