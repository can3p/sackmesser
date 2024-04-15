package simpleyaml

import (
	"gopkg.in/yaml.v3"

	"github.com/can3p/sackmesser/pkg/traverse/simpleobject"
	"github.com/can3p/sackmesser/pkg/traverse/types"
)

type jnode struct {
	types.Node
}

func (n *jnode) Serialize() ([]byte, error) {
	return yaml.Marshal(n.Value())
}

func Parse(b []byte) (types.Node, error) {
	var j any

	if err := yaml.Unmarshal(b, &j); err != nil {
		return nil, err
	}

	return &jnode{
		simpleobject.FromValue(j),
	}, nil
}

func FromNode(n types.Node) types.RootNode {
	return &jnode{
		simpleobject.FromNode(n),
	}
}
