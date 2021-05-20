package d

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Dmap lexer.
var Dmap = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Dmap",
		Aliases:   []string{"dmap"},
		Filenames: []string{"*.dmap"},
		MimeTypes: []string{"text/x-dmap"},
		CaseInsensitive: true,
	},
	dmapRules,
))

// NOTE: The dmap grammar is defined in NX-NASTRAN dmap manual. This lexer is
// preliminary and can only recoganize comments and strings. More complex
// grammars can be added.
func dmapRules() Rules {
	return Rules{
		// Dmap distinguishes lines, so we first chop lines. This is enabled
		// together by requiring multiline regex mode.
		"root": {
			{`^[ ]*\n`, Text, nil},
			{`^\$.*\n`, Comment, nil},
			{`^(.{72})(.+\n)`, ByGroups(UsingSelf("line"), Comment), nil},
			{`^.+\n`, UsingSelf("line"), nil},
		},
		"line": {
			{`([^$]+)(\$.+\n?)`, ByGroups(UsingSelf("statement"), Comment), nil},
			{`([^$]+)(\$\n)`, ByGroups(UsingSelf("statement"), Text), nil},
			{`(.+)(\n?)`, ByGroups(UsingSelf("statement"), Text), nil},
		},
		"statement": {
			{`(.*)('.*?')(.*)`, ByGroups(Text, LiteralString, Text), nil},
			{`.*`, Text, nil},
		},
	}
}
