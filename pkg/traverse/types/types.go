package types

import "fmt"

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

type Node interface {
	Visit(field string) (Node, error)
	NodeType() NodeType
	Value() any
	GetField(field string) (any, error)
	SetField(field string, value any) error
	DeleteField(field string) error
	Serialize() ([]byte, error)
}
