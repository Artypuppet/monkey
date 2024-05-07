package ast

import token "github.com/Artypuppet/monkey/token"

// interface that contains a method that should return the
// token literal with which the node is associated with.
type Node interface {
	// This method is used only for debugging purposes.
	TokenLiteral() string
}

// Dummy interface to help us catch errors in places
// where an expresion should have been used instead of a statement
type Statement interface {
	Node
	statementNode()
}

// Dummy interface to catch errors in places where a statement should have
// been used instead of an expresion.
type Expression interface {
	Node
	expressionNode()
}

// struct representing the root node of the Abstrat Syntax Tree
// Every program consists of a series of statements which are stored
// in Statements slice.
type Program struct {
	Statements []Statement
}

// implementing the Node interface TokenLiteral func returning
// the token literal.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// struct representing the let statement node in the AST
// implements the Statement Interface.
type LetStatement struct {
	Token *token.Token // This is the token.Let token.
	Name  *Identifier  // This is contains the name of the identifier token
	Value *Expression  // this is rhs of the let statement.
}

// empty method to satisfy the Statement interface.
func (ls *LetStatement) statementNode() {}

// implementing the Node Interface.
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// struct defining the identifier in a let statement.
// It implements the expression interface since
// identifier might be produce values in a different statement
// e.g. let x = 5; let y = x; Here in the second statement x is an expression.
type Identifier struct {
	Token *token.Token // This is the token.IDENT token.
	Value string
}

// empty method to implement Expression interface.
func (i *Identifier) expressionNode() {}

// method implementing the Node interface
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
