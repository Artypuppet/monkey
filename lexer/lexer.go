package lexer

import (
	token "github.com/Artypuppet/monkey/token"
)

// TODO change the lexer to support UTF-8 encoding.
// this struct defines the lexer for our language which
// creates tokens character by character
type Lexer struct {
	input        string // The input file/string
	ch           byte   // the current character in input
	position     int    // represents the index of the current ch character in the input
	readPosition int    // represents the index of the next character after ch in the input
}

// constructor for lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// helper method to get the next character in the input string
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() *token.Token {
	var tok *token.Token
	// ignore any whitespace between characters
	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = &token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = &token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case 0:
		tok = newToken(token.EOF, 0)
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			return &token.Token{Type: token.LookupIdent(literal), Literal: literal}
		} else if isDigitFirst(l.ch) {
			return &token.Token{Type: token.INT, Literal: l.readDigit()}
		} else {
			return newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// TODO handle cases where when reading the identifier we encounter a char that is not a letter e.g. %
// helper method to read the an identifier
func (l *Lexer) readIdentifier() string {
	initialPos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[initialPos:l.position]
}

// TODO handle cases where we have a float or right after a digit we have a non digit char
// for floats if we have a number like 9. we interpret it it as a float rather than throwing an error.
// helper function to determine whether
func (l *Lexer) readDigit() string {
	initialPos := l.position
	for isDigitFirst(l.ch) {
		l.readChar()
	}
	return l.input[initialPos:l.position]
}

// helper method to ignore whitespace between characters
func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// helper method to get the next character without advancing the pointers.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// function to create a new token with the given TokenType and Literal value ch and returns a pointer to it.
func newToken(tokenType token.TokenType, ch byte) *token.Token {
	return &token.Token{Type: tokenType, Literal: string(ch)}
}

// function to determine whether the character is a part of the language
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// function that checks if the character is a digit or a plus or a minus sign
// only used for the first character when identifying a number
func isDigitFirst(ch byte) bool {
	return isDigit(ch) || ch == '-' || ch == '+'
}

// function that check if a character is a digit.
func isDigit(ch byte) bool {
	return ('0' <= ch && ch <= '9')
}
