package types

import (
	"bytes"
	"fmt"
	"strconv"
)

var ErrFieldMissing = fmt.Errorf("No field with such name")
var ErrIdxOutOfBounds = fmt.Errorf("Index out of bounds")
var ErrWrongVisit = fmt.Errorf("Not possible to visit a field for the value of this type")

type NodeType string

var (
	NodeTypeString NodeType = "string"
	NodeTypeNumber NodeType = "number"
	NodeTypeBool   NodeType = "bool"
	NodeTypeObject NodeType = "object"
	NodeTypeArray  NodeType = "array"
	NodeTypeNull   NodeType = "null"
)

type PathElement struct {
	ObjectField string
	ArrayIdx    int
}

func (pe PathElement) String() string {
	if pe.ObjectField != "" {
		return pe.ObjectField
	}

	return strconv.Itoa(pe.ArrayIdx)
}

type PathElementSlice []PathElement

func (sl PathElementSlice) String() string {
	var buf bytes.Buffer

	for idx := 0; idx < len(sl); idx++ {
		p := sl[idx]

		if p.ObjectField != "" {
			if idx != 0 {
				buf.WriteRune('.')
			}
			buf.WriteString(p.ObjectField)
			continue
		}
		buf.WriteRune('[')
		buf.WriteString(p.ObjectField)
		buf.WriteString(strconv.Itoa(p.ArrayIdx))
		buf.WriteRune(']')
	}

	return buf.String()
}

type Node interface {
	Visit(field PathElement) (Node, error)
	NodeType() NodeType
	Value() any
	GetField(field PathElement) (any, error)
	SetField(field PathElement, value any) error
	DeleteField(field PathElement) error
}

type RootNode interface {
	Node
	Serialize() ([]byte, error)
}
