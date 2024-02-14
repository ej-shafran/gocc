package lexer

import (
	"bufio"
	"errors"
	"io"
	"unicode"
)

type TokenKind int

const (
	EOF TokenKind = iota + 1
	Ident
	LeftParen
	RightParen
	LeftCurly
	IntegerLiteral
	Semicolon
	RightCurly
)

type Token struct {
	Pos   Position
	Kind  TokenKind
	Value string
}

type Position struct {
	Line int
	Col  int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{Line: 1, Col: 1},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) readRune() (rune, error) {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		return 0, err
	}
	l.pos.Col++
	return r, nil
}

func (l *Lexer) unreadRune() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.Col--
}

func (l *Lexer) nextLine() {
	l.pos.Line++
	l.pos.Col = 1
}

func (l *Lexer) NextToken() (*Token, error) {
	for {
		startPos := l.pos
		r, err := l.readRune()
		if err != nil {
			if err == io.EOF {
				return &Token{Kind: EOF, Pos: l.pos, Value: ""}, nil
			}

			return nil, err
		}

		switch r {
		case '\n':
			l.nextLine()
		case '(':
			return &Token{Pos: startPos, Kind: LeftParen, Value: string(r)}, nil
		case ')':
			return &Token{Pos: startPos, Kind: RightParen, Value: string(r)}, nil
		case '{':
			return &Token{Pos: startPos, Kind: LeftCurly, Value: string(r)}, nil
		case ';':
			return &Token{Pos: startPos, Kind: Semicolon, Value: string(r)}, nil
		case '}':
			return &Token{Pos: startPos, Kind: RightCurly, Value: string(r)}, nil
		default:
			if unicode.IsSpace(r) {
				continue
			}

			if unicode.IsDigit(r) {
				l.unreadRune()
				return l.lexIntegerLiteral()
			}

			if unicode.IsLetter(r) {
				l.unreadRune()
				return l.lexIdentifier()
			}
			return nil, errors.New("TODO: '" + string(r) + "'")
		}
	}
}

func (l *Lexer) lexIdentifier() (token *Token, err error) {
	startPos := l.pos

	var lit []rune
	for {
		r, err := l.readRune()
		if err != nil {
			if err == io.EOF {
				return &Token{Pos: startPos, Value: string(lit), Kind: Ident}, nil
			}

			return nil, err
		}

		if unicode.IsLetter(r) {
			lit = append(lit, r)
		} else {
			l.unreadRune()
			return &Token{Pos: startPos, Value: string(lit), Kind: Ident}, nil
		}
	}
}

func (l *Lexer) lexIntegerLiteral() (token *Token, err error) {
	startPos := l.pos

	var lit []rune
	for {
		r, err := l.readRune()
		if err != nil {
			if err == io.EOF {
				return &Token{Pos: startPos, Value: string(lit), Kind: IntegerLiteral}, nil
			}

			return nil, err

		}

		if unicode.IsDigit(r) {
			lit = append(lit, r)
		} else {
			l.unreadRune()
			return &Token{Pos: startPos, Value: string(lit), Kind: IntegerLiteral}, nil
		}
	}
}
