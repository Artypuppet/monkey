package parser

import (
	"fmt"

	ast "github.com/Artypuppet/monkey/ast"
	lexer "github.com/Artypuppet/monkey/lexer"
	token "github.com/Artypuppet/monkey/token"
)

// struct defining the parser for our program
// holds a lexer to get tokens and curToken
// rep. the current token and peektoken rep.
// the token after curToken
type Parser struct {
	l         *lexer.Lexer
	curToken  *token.Token
	peekToken *token.Token
	errors    []string
}

// Function that creates a new Parser
// takes a lexer for an input as argument
// and call nextToken() twice to set curToken
// and peekToken
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()
	return p
}

// method that moves the curToken and peekToken forward
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

}

// method that parses the program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// parses a statement based on the token type of the cur token.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

// function to parse a let statement.
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.expectPeek(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// checks if current token is of the tokenType
func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

// checks if next token is of the tokenType
func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

// check if token is of the tokenType and accordingly advances the curToken.
func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	} else {
		p.peekError(tokenType)
		return false
	}
}

// getter to return the error slice.
func (p *Parser) Errors() []string {
	return p.errors
}

// method that appends a new error to the errors slice
// for nextToken.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
