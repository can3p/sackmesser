package simpleobject

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/can3p/sackmesser/pkg/traverse/types"
)

type jnode struct {
	v             any
	parent        *jnode
	accessedField types.PathElement
	vType         reflect.Type
}

// if you get a node, it's your responsibility
// to avoid breaking it. As a general rule,
// don't keep the nodes and do not reuse child
// nodes whenever you've done any modifications to
// the parent node
func (n *jnode) Visit(field types.PathElement) (types.Node, error) {
	var err error

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

		if field.ObjectField == "" {
			return nil, types.ErrWrongVisit
		}

		val, ok := m[field.ObjectField]

		if !ok {
			return nil, types.ErrFieldMissing
		}

		return &jnode{
			v:             val,
			parent:        n,
			accessedField: field,
			vType:         reflect.TypeOf(val),
		}, nil
	case types.NodeTypeArray:
		m := n.v.([]any)

		var idx int64

		if field.ObjectField == "" {
			idx = int64(field.ArrayIdx)
		} else {

			idx, err = strconv.ParseInt(field.ObjectField, 10, 64)

			if err != nil {
				return nil, err
			}
		}

		if idx < 0 || idx >= int64(len(m)) {
			return nil, types.ErrIdxOutOfBounds
		}

		val := m[idx]

		return &jnode{
			v:             val,
			parent:        n,
			accessedField: field,
			vType:         reflect.TypeOf(val),
		}, nil
	}

	panic("unreachable")
}

func (n *jnode) GetField(field types.PathElement) (any, error) {
	var err error
	typedArr, ok := n.v.([]any)

	if ok {
		var idx int64

		if field.ObjectField == "" {
			idx = int64(field.ArrayIdx)
		} else {

			idx, err = strconv.ParseInt(field.ObjectField, 10, 64)

			if err != nil {
				return nil, err
			}
		}

		if idx < 0 || idx >= int64(len(typedArr)) {
			return nil, types.ErrIdxOutOfBounds
		}

		return typedArr[idx], nil
	}

	if field.ObjectField == "" {
		return nil, types.ErrWrongVisit
	}

	typed, ok := n.v.(map[string]any)

	if !ok {
		return nil, types.ErrWrongVisit
	}

	if value, ok := typed[field.ObjectField]; ok {
		return value, nil
	}

	return nil, types.ErrFieldMissing
}

func (n *jnode) SetField(field types.PathElement, value any) error {
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
		if field.ObjectField == "" {
			return types.ErrWrongVisit
		}

		m := n.v.(map[string]any)

		m[field.ObjectField] = value

		return nil
	case types.NodeTypeArray:
		m := n.v.([]any)

		var idx int64
		var err error

		if field.ObjectField == "" {
			idx = int64(field.ArrayIdx)
		} else {

			idx, err = strconv.ParseInt(field.ObjectField, 10, 64)

			if err != nil {
				return err
			}
		}

		if idx < 0 || idx >= int64(len(m)) {
			return types.ErrIdxOutOfBounds
		}

		m[idx] = value

		return nil
	}

	panic("unreachable")
}

func (n *jnode) DeleteField(field types.PathElement) error {
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

		if field.ObjectField == "" {
			return types.ErrWrongVisit
		}

		// nothing to do
		if _, ok := m[field.ObjectField]; !ok {
			return nil
		}

		delete(m, field.ObjectField)

		return nil
	case types.NodeTypeArray:
		m := n.v.([]any)

		var idx int64
		var err error

		if field.ObjectField == "" {
			idx = int64(field.ArrayIdx)
		} else {

			idx, err = strconv.ParseInt(field.ObjectField, 10, 64)

			if err != nil {
				return err
			}
		}

		if idx < 0 || idx >= int64(len(m)) {
			return types.ErrIdxOutOfBounds
		}

		// in case of array splice, just mutating the current node
		// value is not enough - parent object somehow retains
		// a reference to the whole slice with initial length.
		// The way to cure that was to modify the array and ask
		// the parent to replace it completely with updated value
		copy(m[idx:], m[idx+1:])
		return n.parent.SetField(n.accessedField, m[:len(m)-1])
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
	// since json does not have them

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

func FromNode(n types.Node) types.Node {
	j := n.Value()

	return &jnode{
		v:     j,
		vType: reflect.TypeOf(j),
	}
}

func FromValue(j any) types.Node {
	return &jnode{
		v:     j,
		vType: reflect.TypeOf(j),
	}
}
