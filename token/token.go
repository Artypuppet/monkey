package token

// type def for string to TokenType to signify the type of
// the token e.g. INT or FUNCTION etc. This helps distinguish
// whether token is a literal, identifier or keyword.
type TokenType string

// struct defining the Tokens in the language
// each token has a TokenType and the value/literal
// of it e.g. in the experession let x = 5
// 'let' is a keyword with Literal value of 'let'
// while 'x' is an IDENTIFIER with a Literal value of 'x' and so on.
type Token struct {
	Type    TokenType
	Literal string
}

// Following are the possible TokenTypes in the language
// ILLEGAL specifies an unknown token type while EOF stand for "end of file"
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	STRING = "STRING" // anything enclosed within ""
	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="
	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// map to emulate a set for faster lookup than switch
var TokenTypes = map[TokenType]struct{}{
	"ILLEGAL": {},
	"":        {}, // Eof
	// Identifiers + literals
	"IDENT": {}, // add, foobar, x, y, ...
	"INT":   {}, // 1343456
	// Operators
	"=": {},
	"+": {},
	// Delimiters
	",": {},
	";": {},
	"(": {},
	")": {},
	"{": {},
	"}": {},
	// Keywords
	"FUNCTION": {},
	"LET":      {},
}

var Idents = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// function that determines whether a string literal is a keyword or an Identifier
func LookupIdent(literal string) TokenType {
	if tokenType, ok := Idents[literal]; ok {
		return tokenType
	}
	return IDENT
}
