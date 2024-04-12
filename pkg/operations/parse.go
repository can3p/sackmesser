package operations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/scanner"

	tscanner "text/scanner"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/can3p/sackmesser/pkg/traverse/types"
	"github.com/pkg/errors"
)

type Operation func(root types.Node, path []string, args ...any) error

type OpInstance struct {
	Op   Operation
	Name string
	Path []string
	Args []any
}

func (op *OpInstance) String() string {
	return fmt.Sprintf("Op: %s, Path: %s, Args: %v", op.Name, strings.Join(op.Path, "."), op.Args)
}

func (op *OpInstance) Apply(root types.Node) error {
	return op.Op(root, op.Path, op.Args...)
}

var operations = map[string]Operation{
	"set": Set,
	"del": Delete,
}

type Call struct {
	Name      string     `@Ident`
	Path      []string   `"(" @Ident ( "." @Ident )*`
	Arguments []Argument `( "," @@ )* ")"`
}

type Argument struct {
	Float  *float64 `  @Float`
	Int    *int     `| @Int`
	Bool   *Boolean `| @("true" | "false")`
	Null   bool     `| @"null"`
	String *String  `| @@`
	JSON   *JSON    `| @@`
}

type String string

// open quote > close quote
var stringQuotes = map[string]string{
	`"`: `"`,
	`'`: `'`,
	"`": "`",
}

var stringEscape = `\`

func (s *String) Parse(lex *lexer.PeekingLexer) error {
	closingQuote, openingFound := stringQuotes[lex.Peek().Value]
	fmt.Println("GOT HERE", lex.RawPeek().Value, openingFound)

	if !openingFound {
		return participle.NextMatch
	}

	var buf bytes.Buffer
	var escaped bool
	var endFound bool

	buf.WriteString(lex.Next().Value)
	fmt.Println(buf.String())

	for {
		token := lex.Next()

		switch {
		case token.EOF():
			return errors.Errorf("EOF reached")
		case escaped:
			escaped = false
		case token.Value == stringEscape:
			escaped = true
		case token.Value == closingQuote:
			endFound = true
		}

		if !escaped {
			buf.WriteString(token.Value)
			fmt.Println(buf.String())
		}

		if endFound {
			break
		}
	}

	*s = String(buf.String())
	return nil
}

type JSON struct {
	Val any
}

// I'm sure we can have a better parser that does not try
// to decode string in a loop, however this would most
// probably require bringing in json parsing in sackmesser
// and I want to keep it light.
//
// In addition, the assumption is that most of the json
// objects in the arguments would be very small
func (j *JSON) Parse(lex *lexer.PeekingLexer) error {
	var buf bytes.Buffer

	token := lex.Peek()

	if token.Value != "[" && token.Value != "{" {
		// unpeek there
		return participle.NextMatch
	}

	buf.WriteString(lex.Next().Value)

	var val any

	for {
		peeked := lex.Next()

		if peeked.EOF() {
			return errors.Errorf("EOF reached")
		}

		buf.WriteString(peeked.Value)

		if err := json.Unmarshal(buf.Bytes(), &val); err == nil {
			j.Val = val
			return nil
		}
	}
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
		participle.Lexer(lexer.NewTextScannerLexer(func(scanner *scanner.Scanner) {
			scanner.Mode = tscanner.ScanIdents | tscanner.ScanFloats | tscanner.ScanInts // we take care of the rest later
		})),
	)

	return &Parser{
		parser: parser,
	}
}

// Parse is very naive at the moment, feel free to replace
// with a proper implementation. Just look how we don't handle
// anything except very simple cases!
// We should be parsing things like
// set(field, "123") // set a string
// set(field, 123) // set a number
// set(field, { a: 1 }) // assign an object to a field
// set(field, "{ a: 1 }") // assign an string to a field
// del(field[0].item) // delete a field
// Problems:
// - Only double quotes are supported for strings which makes passing valid json a pain
// - array indexes are not supported
// - set(field, some spaced value) should be possible
func (p *Parser) Parse(s string) (*OpInstance, error) {
	parsed, err := p.parser.ParseString("", s,
		participle.Trace(os.Stderr),
	)

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

	return &OpInstance{
		Op:   op,
		Name: opName,
		Path: parsed.Path,
		Args: args,
	}, nil
}
