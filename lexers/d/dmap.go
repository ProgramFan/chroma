package d

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Dmap lexer.
var Dmap = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "Dmap",
		Aliases:         []string{"dmap"},
		Filenames:       []string{"*.dmap"},
		MimeTypes:       []string{"text/x-dmap"},
		NotMultiline:    false,
		DotAll:          false,
		CaseInsensitive: true,
	},
	dmapRules,
))

func dmapRules() Rules {
	return Rules{
		// Since the dmap language is a command-driven language, there is no
		// deliminators for functions and instructions and data. So we need just one
		// state. The tricky details here is to arrange the patterns in a correct
		// order so that the most specific pattern goes first.
		"root": {
			Include("whitespace"),
			Include("keywords"),
			Include("literials"),
			{`[a-z0-9]+`, Name, nil},
			Include("operators"),
			{`[(),/]`, Punctuation, nil},
			{`.`, Text, nil},
		},
		"whitespace": {
			{`^\$\n`, Comment, nil},
			{`\$.+\n`, Comment, nil},
			{`\$\n`, Text, nil},
			{`\s+`, Text, nil},
		},
		"keywords": {
			{Words(`\b`, `\b`, "and", "or", "not", "xor", "eqv"), Operator, nil},
			{Words(`\b`, `\b`, "always", "never", "true", "false", "nogo"), KeywordConstant, nil},
			{Words(`\b`, `\b`, "if", "then", "else", "endif", "do", "while", "enddo",
				"subdmap", "call", "return", "exit", "end", "jump", "label"), Keyword, nil},
			{Words(`\b`, `\b`, "type", "dbview", "file", "dbequiv", "dbdelete"), KeywordReserved, nil},
		},
		"literals": {
			{`[+-]?\d+\.(\d+)?([de]\d+)?`, LiteralNumberFloat, nil},
			{`[+-]?\d+`, LiteralNumberInteger, nil},
			{`'.+?'`, LiteralString, nil},
		},
		"operators": {
			{`([*]{2}|[*]|\+|-|/)`, Operator, nil},
			{`(<>|<=|><|>=|<|>|=)`, Operator, nil},
			{`&`, Operator, nil},
		},
	}
}
