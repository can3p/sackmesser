package operations

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/can3p/sackmesser/pkg/operations/lexer"
	"github.com/can3p/sackmesser/pkg/traverse/types"
	"github.com/pkg/errors"
)

type Operation func(root types.Node, path []types.PathElement, args ...any) error

type OpInstance struct {
	Op   Operation
	Name string
	Path types.PathElementSlice
	Args []any
}

func (op *OpInstance) String() string {
	return fmt.Sprintf("Op: %s, Path: %s, Args: %v", op.Name, op.Path.String(), op.Args)
}

func (op *OpInstance) Apply(root types.Node) error {
	return op.Op(root, op.Path, op.Args...)
}

var operations = map[string]Operation{
	"set":   Set,
	"del":   Delete,
	"merge": Merge,
	"pop":   Pop,
	"push":  Push,
}

//nolint:govet
type Call struct {
	Name      string        `@Ident`
	Path      []PathElement `"(" @@+`
	Arguments []Argument    `( "," @@ )* ")"`
}

//nolint:govet
type Argument struct {
	Float  *float64 `  @Float`
	Int    *int     `| @Int`
	Bool   *Boolean `| @("true" | "false")`
	Null   bool     `| @"null"`
	String *string  `| @String | @Ident`
	JSON   *JSON    `| @JSON`
}

//nolint:govet
type PathElement struct {
	// potential foot gun there, I did not want to have a
	// leading dot, but I could not write a grammar rule
	// to exclude it from the first match only, hence
	// I've made it optional
	ObjectField StringPathElement   ` "."? (@String | @Ident)`
	ArrayIdx    ArrIndexPathElement ` | @JSON`
}

type StringPathElement string

func (b *StringPathElement) Capture(values []string) error {
	*b = StringPathElement(strings.Trim(values[0], "\"'`"))
	return nil
}

type ArrIndexPathElement int

var arrayAccessRE = regexp.MustCompile(`\[-?\d+\]`)

// we need to do this because lexer will return text like `[0]` as a single token because of json
func (b *ArrIndexPathElement) Capture(values []string) error {
	if !arrayAccessRE.MatchString(values[0]) {
		return errors.Errorf("Not an array lookup")
	}

	idx, err := strconv.Atoi(strings.Trim(values[0], "[]"))

	if err != nil {
		return err
	}

	*b = ArrIndexPathElement(idx)
	return nil
}

type JSON struct {
	Val any
}

// It's unfortunate that we have parse json one more time
// instead of having the value in the token already
func (b *JSON) Capture(values []string) error {
	var val any

	if err := json.Unmarshal([]byte(values[0]), &val); err != nil {
		return err
	}

	(*b).Val = val
	return nil
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

type Parser struct {
	parser *participle.Parser[Call]
}

func NewParser() *Parser {
	parser := participle.MustBuild[Call](
		participle.Lexer(lexer.NewCustomTextScannerLexer()),
	)

	return &Parser{
		parser: parser,
	}
}

// We should be parsing things like
// set(field, "123") // set a string
// set(field, 123) // set a number
// set(field, { a: 1 }) // assign an object to a field
// set(field, "{ a: 1 }") // assign an string to a field
// del(field.item) // delete a field
// Problems:
// - Only double quotes are supported for strings which makes passing valid json a pain
// - array indexes are not supported
func (p *Parser) Parse(s string) (*OpInstance, error) {
	//parsed, err := p.parser.ParseString("", s,
	//participle.Trace(os.Stderr),
	//)
	parsed, err := p.parser.ParseString("", s)

	if err != nil {
		return nil, err
	}

	opName := strings.ToLower(parsed.Name)

	op, opExists := operations[opName]

	if !opExists {
		return nil, errors.Errorf("Operation [%s] is not supported", opName)
	}

	args := []any{}

	for _, arg := range parsed.Arguments {
		switch {
		case arg.Bool != nil:
			args = append(args, bool(*arg.Bool))
		case arg.Int != nil:
			args = append(args, *arg.Int)
		case arg.Float != nil:
			args = append(args, *arg.Float)
		case arg.String != nil:
			args = append(args, strings.Trim(string(*arg.String), "\"'`"))
		case arg.Null:
			args = append(args, nil)
		case arg.JSON != nil:
			args = append(args, arg.JSON.Val)
		}
	}

	// I've duplicated types to keep parsing data structures
	// and traversal api independent
	path := make([]types.PathElement, 0, len(parsed.Path))
	for _, p := range parsed.Path {
		path = append(path, types.PathElement{
			ObjectField: string(p.ObjectField),
			ArrayIdx:    int(p.ArrayIdx),
		})
	}

	return &OpInstance{
		Op:   op,
		Name: opName,
		Path: path,
		Args: args,
	}, nil
}
