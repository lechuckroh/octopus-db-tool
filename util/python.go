package util

var pythonReservedWords = [...]string{
	"and",
	"as",
	"assert",
	"break",
	"class",
	"continue",
	"def",
	"del",
	"elif",
	"else",
	"except",
	"False",
	"finally",
	"for",
	"from",
	"global",
	"if",
	"import",
	"in",
	"is",
	"lambda",
	"None",
	"nonlocal",
	"not",
	"or",
	"pass",
	"raise",
	"return",
	"True",
	"try",
	"while",
	"with",
	"yield",
}

var pythonReservedWordSet *StringSet

// IsPythonReservedWord returns true if 's' is reserved word in python.
func IsPythonReservedWord(s string) bool {
	if pythonReservedWordSet == nil {
		pythonReservedWordSet = NewStringSet()
		for _, word := range pythonReservedWords {
			pythonReservedWordSet.Add(word)
		}
	}
	return pythonReservedWordSet.Contains(s)
}
