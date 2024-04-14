package simpleyaml

import (
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/can3p/sackmesser/pkg/traverse/types"
)

type jnode struct {
	v     any
	vType reflect.Type
}

func (n *jnode) Serialize() ([]byte, error) {
	return yaml.Marshal(n.v)
}

// if you get a node, it's your responsibility
// to avoid breaking it. As a general rule,
// don't keep the nodes and do not reuse child
// nodes whenever you've done any modifications to
// the parent node
func (n *jnode) Visit(field string) (types.Node, error) {
	switch n.NodeType() {
	case types.NodeTypeNull:
		fallthrough
	case types.NodeTypeNumber:
		fallthrough
	case types.NodeTypeString:
		fallthrough
	case types.NodeTypeBool:
		return nil, types.ErrWrongVisit
	case types.NodeTypeObject:
		m := n.v.(map[string]any)

		val, ok := m[field]

		if !ok {
			return nil, types.ErrFieldMissing
		}

		return &jnode{
			v:     val,
			vType: reflect.TypeOf(val),
		}, nil
	case types.NodeTypeArray:
		m := n.v.([]any)

		idx, err := strconv.ParseInt(field, 10, 64)

		if err != nil {
			return nil, err
		}

		if idx < 0 || idx >= int64(len(m)) {
			return nil, types.ErrIdxOutOfBounds
		}

		val := m[idx]

		return &jnode{
			v:     val,
			vType: reflect.TypeOf(val),
		}, nil
	}

	panic("unreachable")
}

func (n *jnode) SetField(field string, value any) error {
	switch n.NodeType() {
	case types.NodeTypeNull:
		fallthrough
	case types.NodeTypeNumber:
		fallthrough
	case types.NodeTypeString:
		fallthrough
	case types.NodeTypeBool:
		return types.ErrWrongVisit
	case types.NodeTypeObject:
		m := n.v.(map[string]any)

		m[field] = value

		return nil
	case types.NodeTypeArray:
		m := n.v.([]any)

		idx, err := strconv.ParseInt(field, 10, 64)

		if err != nil {
			return err
		}

		if idx < 0 || idx >= int64(len(m)) {
			return types.ErrIdxOutOfBounds
		}

		m[idx] = value

		return nil
	}

	panic("unreachable")
}

func (n *jnode) DeleteField(field string) error {
	switch n.NodeType() {
	case types.NodeTypeNull:
		fallthrough
	case types.NodeTypeNumber:
		fallthrough
	case types.NodeTypeString:
		fallthrough
	case types.NodeTypeBool:
		return types.ErrWrongVisit
	case types.NodeTypeObject:
		m := n.v.(map[string]any)

		if _, ok := m[field]; !ok {
			return types.ErrFieldMissing
		}

		delete(m, field)

		return nil
	case types.NodeTypeArray:
		m := n.v.([]any)

		idx, err := strconv.ParseInt(field, 10, 64)

		if err != nil {
			return err
		}

		if idx < 0 || idx >= int64(len(m)) {
			return types.ErrIdxOutOfBounds
		}

		copy(m[idx:], m[idx+1:])
		n.v = m[:len(m)-1]

		return nil
	}

	panic("unreachable")
}

func (n *jnode) Value() any {
	return n.v
}

func (n *jnode) NodeType() types.NodeType {
	// this is a mega ugly hack,
	// see https://github.com/golang/go/issues/51649

	isNil := fmt.Sprintf("%v", n.vType) == "<nil>"
	// we do not enumerate all values of kind
	// since yaml does not have them

	if isNil {
		return types.NodeTypeNull
	}

	switch n.vType.Kind() {
	case reflect.Bool:
		return types.NodeTypeBool
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uintptr:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		return types.NodeTypeNumber
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		return types.NodeTypeArray
	case reflect.Map:
		return types.NodeTypeObject
	case reflect.String:
		return types.NodeTypeString
	}

	panic("Unreachable")
}

func Parse(b []byte) (types.Node, error) {
	var j any

	if err := yaml.Unmarshal(b, &j); err != nil {
		return nil, err
	}

	return &jnode{
		v:     j,
		vType: reflect.TypeOf(j),
	}, nil
}

func FromNode(n types.Node) types.Node {
	j := n.Value()

	return &jnode{
		v:     j,
		vType: reflect.TypeOf(j),
	}
}
