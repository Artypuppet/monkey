package ast

import (
	"bytes"

	token "github.com/Artypuppet/monkey/token"
)

// -------------------Interfaces for Nodes in AST--------------------------

// interface that contains a method that should return the
// token literal with which the node is associated with.
type Node interface {
	// These methods are used only for debugging purposes.
	TokenLiteral() string
	String() string
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

// -------------------------Program struct(Root node)----------------------
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

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// -------------------Let Statement struct for Let Nodes-----------------------

// struct representing the let statement node in the AST
// implements the Statement Interface.
type LetStatement struct {
	Token *token.Token // This is the token.Let token.
	Name  *Identifier  // This is contains the name of the identifier token
	Value Expression   // this is rhs of the let statement.
}

// empty method to satisfy the Statement interface.
func (ls *LetStatement) statementNode() {}

// implementing the Node Interface.
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

// ----------------------------------Identifier--------------------------------

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

func (i *Identifier) String() string {
	return i.Value
}

//--------------------------------Return Statement------------------------------

// struct defining the node associated with a return statement
type ReturnStatement struct {
	Token       *token.Token // token type will be RETURN
	ReturnValue Expression   // This will be an expression e.g. return add(5, 6)
}

// satisfying the statement interface
func (rs *ReturnStatement) statementNode() {}

// satisfying the Node interface
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// --------------------------Expression Statement-------------------------------

// struct defining expression statement
// This is distinct from a simple expression or a statment as
// monkey support the following expreesion statement
// let x = 5;
// x + 10;
type ExpressionStatement struct {
	Token      *token.Token // The first token of the expression e.g. x in x + 10
	Expression Expression
}

// method to satisfy Statement interface
func (es *ExpressionStatement) statementNode() {}

// method to satisfy expression interface
// func (es *ExpressionStatement) expressionNode() {}

// method to satisfy Node interface
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// -----------------------------------Integer Literals---------------------------
// struct that represents an integer literal
// It implements the Expression Interface.
type IntegerLiteral struct {
	Token *token.Token
	Value int64 // The parsed value of Token.Literal
}

// methods to satisfy the Expression Interface
func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// -----------------------------Prefix Expression Node--------------------------
// struct representing a prefix expression
// It implements the Expression Interface
// Operator represents any operator applied to the expression on the right.
// prefix operators consist of --, ++, -, !, +
type PrefixExpression struct {
	Token    *token.Token
	Operator string
	Right    Expression
}

// methods to implement the Expression interface
func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// ------------------------------Infix Expression Node-----------------------------
// struct represents an infix Expression node in the ast
// It implements the Expression interface
// Operators represents the actual infix operator
// Left and right represent the operands of the infix operator
type InfixExpression struct {
	Token    *token.Token // Operator Token. Can be +, -, *, /, <, >, !=, ==
	Operator string
	Left     Expression
	Right    Expression
}

// methods to implement Expression interface
func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
