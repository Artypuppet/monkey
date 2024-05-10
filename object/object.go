package object

import "fmt"

// ------------------------------Object---------------------------------

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
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
///object.Integer, saving the value inside our struct and passing around a reference to
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
