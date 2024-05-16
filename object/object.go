package object

import (
	"bytes"
	"fmt"
	"strings"

	ast "github.com/Artypuppet/monkey/ast"
)

// ------------------------------Object---------------------------------

type ObjectType string

// This is typedef for builtin functions that can be called within monkey
type BuiltinFunction func(args ...Object) Object

const (
	INTEGER_OBJ      = "INTEGER"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
)

// this interface defines the top level value representation of
// Literal values encountered when evaulating AST.
// The reason for opting for an interface is because the values
// could be different it would be difficult to cram all types into
// a single object.
type Object interface {
	Type() ObjectType // returns the type of the object
	Inspect() string  // returns the value, formatted as a string.
}

// ----------------------------Integer Literal---------------------------

// struct defining the internal representation for an integer literal
// It implements the Object interface.
// Whenever we encounter an integer literal in the source code we first turn it into an
// ast.IntegerLiteral and then, when evaluating that AST node, we turn it into an
// object.Integer, saving the value inside our struct and passing around a reference to
// this struct.
type Integer struct {
	Value int64
}

// Methods implementing the Object interface.
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// ----------------------------String Literal-----------------------------

// struct defining the internal representation for a string literal
// It implements the Object interface.
type String struct {
	Value string
}

// Methods implementing the Object interface.
func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

// ----------------------------Boolean Literal----------------------------

// struct defining the internal representation for a boolean literal
// It implements the Object interface.
type Boolean struct {
	Value bool
}

// methods to implement the Object interface
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// ---------------------------Null Literal-------------------------------

// struct defining the internal representation for null.
// It implements the Object interface.
// It's empty because we already know the value for it.
type Null struct{}

// method to implement the Object interface.
func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

// --------------------------Return Value-------------------------------

// struct defining the internal representation for return
// statements. Execution for the current block of code should
// stop when return is ecountered i.e. all subsequent statements
// must be skipped.
// Implements the Object interface.
type ReturnValue struct {
	Value Object
}

// methods to implement the object interface.
func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// --------------------------Error------------------------------------

// struct defining error struct to represent any error that was encountered
// while evaluating the code.
// Implements the object interface
type Error struct {
	Message string
}

// methods to implement the object interface.
func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

// --------------------------------Environment------------------------

// struct defining the environment object to keep track of variables
// bindings etc.
// PITFALL: maps are reference types as in they are not copied when passing
// to a function or returned from a function.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Function to that return an an instance of the Environment struct
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// Function that returns a new Environment with a ptr to its outer environment
func NewEnclosedEnvironment(outer *Environment) *Environment {
	enclosed := NewEnvironment()
	enclosed.outer = outer
	return enclosed
}

// method to get the object associated with an identifier
// The identifier here is node.Name.Value where node is a LetStatement
// Name is an identifier struct and Value is a string.
// It checks if the identifier already exists in the environment
// If not it then calls get in its outer environment.
// This mechanism helps enforce function scope wherein functions
// can reference variables outside then defintion but outer variable with
// the same name as the function parameter will always reference the parameter.
func (e *Environment) Get(identifier string) (Object, bool) {
	obj, ok := e.store[identifier]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(identifier)
	}
	return obj, ok
}

// method to set the object for an identifier
// The identifier here is node.Name.Value where node is a LetStatement
// Name is an identifier struct and Value is a string.
func (e *Environment) Set(identifier string, val Object) Object {
	e.store[identifier] = val
	return val
}

// -----------------------------Function Object-----------------------

// struct to represent function object in our environment
// It implements the Object interface.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// methods to implement the object interface
func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// --------------------------Builtin Function-------------------------

// stuct that is a wrapper around a builtin function
// It implements the object interface.
type Builtin struct {
	Fn BuiltinFunction
}

// methods implementing the object interface.
func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}
