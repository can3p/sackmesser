package operations

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/can3p/sackmesser/pkg/scanner"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/pkg/errors"
)

// butchered version of https://github.com/alecthomas/participle/blob/master/lexer/text_scanner.go

// TextScannerLexer is a lexer that uses the text/scanner module.
var (
	TextScannerLexer lexer.Definition = &textScannerLexerDefinition{}

	// DefaultDefinition defines properties for the default lexer.
	DefaultDefinition = TextScannerLexer
)

// NewCustomTextScannerLexer constructs a Definition that uses an underlying scanner.Scanner
//
// It's custom because:
// - string token can have different sets of quotes: single, double, backtick
// - special token type - JSON
func NewCustomTextScannerLexer() lexer.Definition {
	return &textScannerLexerDefinition{}
}

type textScannerLexerDefinition struct{}

func (d *textScannerLexerDefinition) Lex(filename string, r io.Reader) (lexer.Lexer, error) {
	l := Lex(filename, r)
	return l, nil
}

func (d *textScannerLexerDefinition) Symbols() map[string]lexer.TokenType {
	return map[string]lexer.TokenType{
		"EOF":    lexer.EOF,
		"Ident":  scanner.Ident,
		"Int":    scanner.Int,
		"Float":  scanner.Float,
		"String": scanner.String,
		"JSON":   scanner.JSON,
	}
}

// textScannerLexer is a Lexer based on text/scanner.Scanner
type textScannerLexer struct {
	scanner  *scanner.Scanner
	filename string
	err      error
}

// Lex an io.Reader with text/scanner.Scanner.
//
// This provides very fast lexing of source code compatible with Go tokens.
//
// Note that this differs from text/scanner.Scanner in that string tokens will be unquoted.
func Lex(filename string, r io.Reader) lexer.Lexer {
	s := &scanner.Scanner{}
	s.Init(r)
	lexerStruct := lexWithScanner(filename, s)
	lexerStruct.scanner.Error = func(s *scanner.Scanner, msg string) {
		lexerStruct.err = &lexer.Error{Msg: msg, Pos: lexer.Position(lexerStruct.scanner.Pos())}
	}
	return lexerStruct
}

// LexWithScanner creates a Lexer from a user-provided scanner.Scanner.
//
// Useful if you need to customise the Scanner.
func LexWithScanner(filename string, scan *scanner.Scanner) lexer.Lexer {
	return lexWithScanner(filename, scan)
}

func lexWithScanner(filename string, scan *scanner.Scanner) *textScannerLexer {
	scan.Filename = filename
	lexer := &textScannerLexer{
		filename: filename,
		scanner:  scan,
	}
	return lexer
}

// LexBytes returns a new default lexer over bytes.
func LexBytes(filename string, b []byte) lexer.Lexer {
	return Lex(filename, bytes.NewReader(b))
}

// LexString returns a new default lexer over a string.
func LexString(filename, s string) lexer.Lexer {
	return Lex(filename, strings.NewReader(s))
}

// open quote > close quote
var stringQuotes = map[rune]rune{
	'"':  '"',
	'\'': '\'',
	'`':  '`',
}

var stringEscape = '\\'

func (t *textScannerLexer) parseString() (string, error) {
	ch := t.scanner.Peek()

	closingQuote, openingFound := stringQuotes[ch]

	if !openingFound {
		return "", errors.Errorf("Unknown string quote: [%c]", ch)
	}

	var buf bytes.Buffer
	var escaped bool
	var endFound bool

	buf.WriteRune(ch)

	for {
		_ = t.scanner.Next()
		ch = t.scanner.Peek()

		switch {
		case ch == '\n' || ch < 0:
			return "", errors.Errorf("Liternal not terminated")
		case escaped:
			escaped = false
		case ch == stringEscape:
			escaped = true
		case ch == closingQuote:
			endFound = true
		}

		if !escaped {
			buf.WriteRune(ch)
		}

		if endFound {
			break
		}
	}

	// advance scanner by one to point to the first non parsed char
	_ = t.scanner.Next()

	return buf.String(), nil
}

func (t *textScannerLexer) parseJSON() (string, error) {
	ch := t.scanner.Peek()
	var buf bytes.Buffer

	if ch != '[' && ch != '{' {
		return "", errors.Errorf("Unexpected start of json input")
	}

	buf.WriteRune(ch)

	var val any

	for {
		_ = t.scanner.Next()
		ch = t.scanner.Peek()

		if ch < 0 {
			return "", errors.Errorf("EOF reached")
		}

		buf.WriteRune(ch)

		if err := json.Unmarshal(buf.Bytes(), &val); err == nil {
			// advance scanner by one to point to the first non parsed char
			_ = t.scanner.Next()
			return buf.String(), nil
		}
	}
}

func (t *textScannerLexer) Next() (lexer.Token, error) {
	typ := t.scanner.Scan()
	text := t.scanner.TokenText()

	pos := lexer.Position(t.scanner.Position)
	pos.Filename = t.filename
	if t.err != nil {
		return lexer.Token{}, t.err
	}
	return lexer.Token{
		Type:  lexer.TokenType(typ),
		Value: text,
		Pos:   pos,
	}, nil
}
