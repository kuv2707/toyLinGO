package globals

var TOKENS=map[string]string{
	"SCOPE_START":"{",
	"SCOPE_END":"}",
	"COLON":":",
	"SEMICOLON":";",
	"LET":"let",
	"INTEGER":"int",
	"FLOAT":"float",
	"BOOLEAN":"bool",
	"STRING":"string",
	"DOT":".",
	"OPEN_PAREN":"(",
	"CLOSE_PAREN":")",
	"IF":"if",
	"ELSE IF":"elseif",
	"ELSE":"else",
	"LOOP":"loop",
	"BREAK":"break",
	"FUNCTION":"func",
	"COMMA":",",
	"RETURN":"return",
}

var KEYWORDS = []string{"IF", "ELSE IF", "ELSE", "FUNCTION", "SCOPE_END", "LET"}