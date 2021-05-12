package o

import (
    . "github.com/alecthomas/chroma" // nolint
    "github.com/alecthomas/chroma/lexers/internal"
)

// Omg Interface Definition Language lexer.
var OmgIdlLexer = internal.Register(MustNewLazyLexer(
    &Config{
        Name:      "OmgIdl",
        Aliases:   []string{ "omg-idl",  },
        Filenames: []string{ "*.idl", "*.pidl",  },
        MimeTypes: []string{  },
    },
    func() Rules {
        return Rules{
            "values": {
                { Words(`(?i)`, `\b`, `true`, `false`), LiteralNumber, nil },
                { `([Ll]?)(")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Push("string") },
                { `([Ll]?)(\')(\\[^\']+)(\')`, ByGroups(LiteralStringAffix, LiteralStringChar, LiteralStringEscape, LiteralStringChar), nil },
                { `([Ll]?)(\')(\\\')(\')`, ByGroups(LiteralStringAffix, LiteralStringChar, LiteralStringEscape, LiteralStringChar), nil },
                { `([Ll]?)(\'.\')`, ByGroups(LiteralStringAffix, LiteralStringChar), nil },
                { `[+-]?\d+(\.\d*)?[Ee][+-]?\d+`, LiteralNumberFloat, nil },
                { `[+-]?(\d+\.\d*)|(\d*\.\d+)([Ee][+-]?\d+)?`, LiteralNumberFloat, nil },
                { `(?i)[+-]?0x[0-9a-f]+`, LiteralNumberHex, nil },
                { `[+-]?[1-9]\d*`, LiteralNumberInteger, nil },
                { `[+-]?0[0-7]*`, LiteralNumberOct, nil },
                { `[\+\-\*\/%^&\|~]`, Operator, nil },
                { Words(``, ``, `<<`, `>>`), Operator, nil },
                { `((::)?\w+)+`, Name, nil },
                { `[{};:,<>\[\]]`, Punctuation, nil },
            },
            "annotation_params": {
                Include("whitespace"),
                { `\(`, Punctuation, Push() },
                Include("values"),
                { `=`, Punctuation, nil },
                { `\)`, Punctuation, Pop(1) },
            },
            "annotation_params_maybe": {
                { `\(`, Punctuation, Push("annotation_params") },
                Include("whitespace"),
                Default(Pop(1)),
            },
            "annotation_appl": {
                { `@((::)?\w+)+`, NameDecorator, Push("annotation_params_maybe") },
            },
            "enum": {
                Include("whitespace"),
                { `[{,]`, Punctuation, nil },
                { `\w+`, NameConstant, nil },
                Include("annotation_appl"),
                { `\}`, Punctuation, Pop(1) },
            },
            "root": {
                Include("whitespace"),
                { Words(`(?i)`, `\b`, `typedef`, `const`, `in`, `out`, `inout`, `local`), KeywordDeclaration, nil },
                { Words(`(?i)`, `\b`, `void`, `any`, `native`, `bitfield`, `unsigned`, `boolean`, `char`, `wchar`, `octet`, `short`, `long`, `int8`, `uint8`, `int16`, `int32`, `int64`, `uint16`, `uint32`, `uint64`, `float`, `double`, `fixed`, `sequence`, `string`, `wstring`, `map`), KeywordType, nil },
                { Words(`(?i)`, `(\s+)(\w+)`, `@annotation`, `struct`, `union`, `bitset`, `interface`, `exception`, `valuetype`, `eventtype`, `component`), ByGroups(Keyword, TextWhitespace, NameClass), nil },
                { Words(`(?i)`, `\b`, `abstract`, `alias`, `attribute`, `case`, `connector`, `consumes`, `context`, `custom`, `default`, `emits`, `factory`, `finder`, `getraises`, `home`, `import`, `manages`, `mirrorport`, `multiple`, `Object`, `oneway`, `primarykey`, `private`, `port`, `porttype`, `provides`, `public`, `publishes`, `raises`, `readonly`, `setraises`, `supports`, `switch`, `truncatable`, `typeid`, `typename`, `typeprefix`, `uses`, `ValueBase`), Keyword, nil },
                { `(?i)(enum|bitmask)(\s+)(\w+)`, ByGroups(Keyword, TextWhitespace, NameClass), Push("enum") },
                { `(?i)(module)(\s+)(\w+)`, ByGroups(KeywordNamespace, TextWhitespace, NameNamespace), nil },
                { `(\w+)(\s*)(=)`, ByGroups(NameConstant, TextWhitespace, Operator), nil },
                { `[\(\)]`, Punctuation, nil },
                Include("values"),
                Include("annotation_appl"),
            },
            "keywords": {
                { Words(``, `\b`, `_Alignas`, `_Alignof`, `_Noreturn`, `_Generic`, `_Thread_local`, `_Static_assert`, `_Imaginary`, `noreturn`, `imaginary`, `complex`), Keyword, nil },
                { `(struct|union)(\s+)`, ByGroups(Keyword, Text), Push("classname") },
                { Words(``, `\b`, `asm`, `auto`, `break`, `case`, `const`, `continue`, `default`, `do`, `else`, `enum`, `extern`, `for`, `goto`, `if`, `register`, `restricted`, `return`, `sizeof`, `struct`, `static`, `switch`, `typedef`, `volatile`, `while`, `union`, `thread_local`, `alignas`, `alignof`, `static_assert`, `_Pragma`), Keyword, nil },
                { Words(``, `\b`, `inline`, `_inline`, `__inline`, `naked`, `restrict`, `thread`), KeywordReserved, nil },
                { `(__m(128i|128d|128|64))\b`, KeywordReserved, nil },
                { Words(`__`, `\b`, `asm`, `based`, `except`, `stdcall`, `cdecl`, `fastcall`, `declspec`, `finally`, `try`, `leave`, `w64`, `unaligned`, `raise`, `noop`, `identifier`, `forceinline`, `assume`), KeywordReserved, nil },
            },
            "types": {
                { Words(``, `\b`, `_Bool`, `_Complex`, `_Atomic`), KeywordType, nil },
                { Words(`__`, `\b`, `int8`, `int16`, `int32`, `int64`, `wchar_t`), KeywordReserved, nil },
                { Words(``, `\b`, `bool`, `int`, `long`, `float`, `short`, `double`, `char`, `unsigned`, `signed`, `void`), KeywordType, nil },
            },
            "whitespace": {
                { `^#if\s+0`, CommentPreproc, Push("if0") },
                { `^#`, CommentPreproc, Push("macro") },
                { `^(\s*(?:/[*].*?[*]/\s*)?)(#if\s+0)`, ByGroups(UsingSelf("root"), CommentPreproc), Push("if0") },
                { `^(\s*(?:/[*].*?[*]/\s*)?)(#)`, ByGroups(UsingSelf("root"), CommentPreproc), Push("macro") },
                { `\n`, Text, nil },
                { `\s+`, Text, nil },
                { `\\\n`, Text, nil },
                { `//(\n|[\w\W]*?[^\\]\n)`, CommentSingle, nil },
                { `/(\\\n)?[*][\w\W]*?[*](\\\n)?/`, CommentMultiline, nil },
                { `/(\\\n)?[*][\w\W]*`, CommentMultiline, nil },
            },
            "statements": {
                Include("keywords"),
                Include("types"),
                { `([LuU]|u8)?(")`, ByGroups(LiteralStringAffix, LiteralString), Push("string") },
                { `([LuU]|u8)?(')(\\.|\\[0-7]{1,3}|\\x[a-fA-F0-9]{1,2}|[^\\\'\n])(')`, ByGroups(LiteralStringAffix, LiteralStringChar, LiteralStringChar, LiteralStringChar), nil },
                { `0[xX]([0-9a-fA-F](\'?[0-9a-fA-F])*\.[0-9a-fA-F](\'?[0-9a-fA-F])*|\.[0-9a-fA-F](\'?[0-9a-fA-F])*|[0-9a-fA-F](\'?[0-9a-fA-F])*)[pP][+-]?[0-9a-fA-F](\'?[0-9a-fA-F])*[lL]?`, LiteralNumberFloat, nil },
                { `(-)?(\d(\'?\d)*\.\d(\'?\d)*|\.\d(\'?\d)*|\d(\'?\d)*)[eE][+-]?\d(\'?\d)*[fFlL]?`, LiteralNumberFloat, nil },
                { `(-)?((\d(\'?\d)*\.(\d(\'?\d)*)?|\.\d(\'?\d)*)[fFlL]?)|(\d(\'?\d)*[fFlL])`, LiteralNumberFloat, nil },
                { `(-)?0[xX][0-9a-fA-F](\'?[0-9a-fA-F])*(([uU][lL]{0,2})|[lL]{1,2}[uU]?)?`, LiteralNumberHex, nil },
                { `(-)?0[bB][01](\'?[01])*(([uU][lL]{0,2})|[lL]{1,2}[uU]?)?`, LiteralNumberBin, nil },
                { `(-)?0(\'?[0-7])+(([uU][lL]{0,2})|[lL]{1,2}[uU]?)?`, LiteralNumberOct, nil },
                { `(-)?\d(\'?\d)*(([uU][lL]{0,2})|[lL]{1,2}[uU]?)?`, LiteralNumberInteger, nil },
                { `[~!%^&*+=|?:<>/-]`, Operator, nil },
                { `[()\[\],.]`, Punctuation, nil },
                { `(true|false|NULL)\b`, NameBuiltin, nil },
                { `((?:[a-zA-Z_$]|\\u[0-9a-fA-F]{4}|\\U[0-9a-fA-F]{8})(?:[\w$]|\\u[0-9a-fA-F]{4}|\\U[0-9a-fA-F]{8})*)(\s*)(:)(?!:)`, ByGroups(NameLabel, Text, Punctuation), nil },
                { `(?:[a-zA-Z_$]|\\u[0-9a-fA-F]{4}|\\U[0-9a-fA-F]{8})(?:[\w$]|\\u[0-9a-fA-F]{4}|\\U[0-9a-fA-F]{8})*`, Name, nil },
            },
            "statement": {
                Include("whitespace"),
                Include("statements"),
                { `\}`, Punctuation, nil },
                { `[{;]`, Punctuation, Pop(1) },
            },
            "function": {
                Include("whitespace"),
                Include("statements"),
                { `;`, Punctuation, nil },
                { `\{`, Punctuation, Push() },
                { `\}`, Punctuation, Pop(1) },
            },
            "string": {
                { `"`, LiteralString, Pop(1) },
                { `\\([\\abfnrtv"\']|x[a-fA-F0-9]{2,4}|u[a-fA-F0-9]{4}|U[a-fA-F0-9]{8}|[0-7]{1,3})`, LiteralStringEscape, nil },
                { `[^\\"\n]+`, LiteralString, nil },
                { `\\\n`, LiteralString, nil },
                { `\\`, LiteralString, nil },
            },
            "macro": {
                { `(\s*(?:/[*].*?[*]/\s*)?)(include)(\s*(?:/[*].*?[*]/\s*)?)("[^"]+")([^\n]*)`, ByGroups(UsingSelf("root"), CommentPreproc, UsingSelf("root"), CommentPreprocFile, CommentSingle), nil },
                { `(\s*(?:/[*].*?[*]/\s*)?)(include)(\s*(?:/[*].*?[*]/\s*)?)(<[^>]+>)([^\n]*)`, ByGroups(UsingSelf("root"), CommentPreproc, UsingSelf("root"), CommentPreprocFile, CommentSingle), nil },
                { `[^/\n]+`, CommentPreproc, nil },
                { `/[*](.|\n)*?[*]/`, CommentMultiline, nil },
                { `//.*?\n`, CommentSingle, Pop(1) },
                { `/`, CommentPreproc, nil },
                { `(?<=\\)\n`, CommentPreproc, nil },
                { `\n`, CommentPreproc, Pop(1) },
            },
            "if0": {
                { `^\s*#if.*?(?<!\\)\n`, CommentPreproc, Push() },
                { `^\s*#el(?:se|if).*\n`, CommentPreproc, Pop(1) },
                { `^\s*#endif.*?(?<!\\)\n`, CommentPreproc, Pop(1) },
                { `.*?\n`, Comment, nil },
            },
            "classname": {
                { `(?:[a-zA-Z_$]|\\u[0-9a-fA-F]{4}|\\U[0-9a-fA-F]{8})(?:[\w$]|\\u[0-9a-fA-F]{4}|\\U[0-9a-fA-F]{8})*`, NameClass, Pop(1) },
                { `\s*(?=>)`, Text, Pop(1) },
                Default(Pop(1)),
            },
        }
    },
))

