package parser

import (
	"fmt"
	"strconv"

	ast "github.com/Artypuppet/monkey/ast"
	lexer "github.com/Artypuppet/monkey/lexer"
	token "github.com/Artypuppet/monkey/token"
)

// ------------------------------------Parser-----------------------------------

// Parse functions for infix and prefix tokens/expressions
// for Pratt's Parsing
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
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
	// maps for tokenTypes and their associated parse functions.
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// Function that creates a new Parser
// takes a lexer for an input as argument
// and call nextToken() twice to set curToken
// and peekToken
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()

	// Register prefix parse fns.
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	// Register infix parse fns.
	// All token types here are associated with the same function
	// as it can handle all of them.
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	return p
}

// helper methods to register parse functions
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
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

// -----------------------------------Helper methods to parse statements----------------
// parses a statement based on the token type of the cur token.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
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

// method that appends a new error to the errors slice for nextToken.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// ------------------------------Let Statement Parsing---------------------------------

// function to parse a let statement.
// It creates a LetStatement struct
// and then checks if it is followed by an identifier
// and if the identifier is followed by an assignment sign
// and subsequently parses the expression after it.
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// curToken is still '=' so we move the token forward.
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	// we move the curToken forward if we encounter a ';' since they are optional so they don't add anything.
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// ---------------------------Return Statement Parsing---------------------------------

// function to parse a return statement
// it checks if there is an expression after Return
// and if there is one it parses it.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	// move to the next token
	p.nextToken()

	// parse the return Expression.
	stmt.ReturnValue = p.parseExpression(LOWEST)

	// advance the curToken if the next token is semicolon.
	p.expectPeek(token.SEMICOLON)

	return stmt
}

// -----------------------------Parse Expression Statement----------------------------

// parsing precedence as an enum essentially.
// _ in the first line means 0, therefore all others are given increasing values
// from 1-7.
const (
	_           int = iota
	LOWEST          // 1
	EQUALS          // 2 ==
	LESSGREATER     // 3 > or <
	SUM             // 4 +
	PRODUCT         // 5 *
	PREFIX          // 6 -X or !X
	CALL            // 7 myFunction(X)
	INDEX           // 8 []
)

// map that defines precedences of different token types
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

// helper method to check the precendence of the next Token
// While parsing infix expression we first encounter the left expression
// we peek to the next Token to get its precedence.
// It returns the LOWEST precedence for any token that does not exist in the map.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// helper method to check the precedence of the current token.
// If the token is not an operators than it certainly will have the lowest precedence.
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// function to parse Expression Statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// method to handle prefix parsing errors.
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// method that parses expressions.
// the prefix could be an operator like ! or -, or it could be an
// identifier or an integer.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := p.prefixParseFns[p.curToken.Type]
	if prefixFn == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefixFn()
	for !p.curTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixFn := p.infixParseFns[p.peekToken.Type]
		if infixFn == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infixFn(leftExp)
	}
	return leftExp
}

// ----------------------------Parse Identifier------------------------------------
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// -------------------------Parse Integer Literal----------------------------------
func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: val}
}

// -------------------------Parse String Literal----------------------------------
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// -------------------------Parse Prefix Expression--------------------------------
// This function parses the operator and then calls the parseExpression method to
// parse the rest of the expression.
func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}

	p.nextToken()

	exp.Right = p.parseExpression(PREFIX)

	return exp
}

// -------------------------Parse Infix Expression--------------------------------

// This function parses infix expression by taking in the left operand for the operator
// and then calling parseExpression() for the right hand operand with the precedence of the
// operator.
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{Token: p.curToken, Operator: p.curToken.Literal, Left: left}

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}

// --------------------------Parse Boolean Expression-----------------------------

// Function to parse boolean expressions
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// -------------------------Parse Grouped Expression------------------------------

// function to parse grouped expressions
// it Reads the left parentheses and then calls parse expression to parse expression
// within the parentheses and then checks if the parentheses are balanced in the end.
func (p *Parser) parseGroupedExpression() ast.Expression {

	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// ----------------------------------Parse If Expression---------------------------

// function to parse If Expression
func (p *Parser) parseIfExpression() ast.Expression {

	exp := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}
	return exp
}

// -------------------------------Parse Block Statement------------------------------

// Function to parse Block statements
// block statements are special because the last statement in the block
// is the value that is returned from the block e.g. let exp = if (x > 5) { 7 } else { 8 }
// in this case exp will be assigned either 7 or 8
func (p *Parser) parseBlockStatement() *ast.BlockStatement {

	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func trace(fnName string) string {
	fmt.Printf("BEGIN %s\n", fnName)
	return fnName
}

func untrace(fnName string) {
	fmt.Printf("END %s\n", fnName)
}

// -------------------------Parse Function Literal Expression-------------------------
// function that parses a function literal
// The current toke when this function called is 'fn'
func (p *Parser) parseFunctionLiteral() ast.Expression {

	lit := &ast.FunctionLiteral{Token: p.curToken}

	// move forward the token if expectPeek is true curToken is IF at this point
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// curToken is now '('
	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {

	parameters := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}

	p.nextToken()

	// curToken is now some identifier
	identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	parameters = append(parameters, identifier)

	for p.peekTokenIs(token.COMMA) {
		// move the curToken such that we skip the comma and then curToken is again some identifier
		p.nextToken()
		p.nextToken()
		parameters = append(parameters, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return parameters
}

// -------------------------------------Parse Call Expression------------------------

// this function is called whenever a '(' is encountered when in the parse Expression
// function. This means that it will be called after an identifier has been parsed or
// it will be called after '}'
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	// when this fn is called the curToken is '('
	callExp := &ast.CallExpression{Token: p.curToken, Function: function}
	callExp.Arguments = p.parseExpressionList(token.RPAREN)
	return callExp
}

// function that parses the function call arguments
// curToken is still '(' when this function is called.
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	// if the next token is ')' then this fn has no arguments.
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

// This function generalizes instances where we need to parse a list.
// e.g. function arguments or arrays.
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}

// This function parses the index expression for arrays.
// It is called for an infix expression. The expression
// between [] should produce an integer.
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}
