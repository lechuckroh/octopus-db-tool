package util

import (
	"github.com/iancoleman/strcase"
	"github.com/xwb1989/sqlparser"
	"strconv"
	"strings"
	"text/template"
)

func ToInt(value interface{}, defaultValue int) int {
	switch value.(type) {
	case string:
		i, err := strconv.Atoi(value.(string))
		if err != nil {
			return defaultValue
		}
		return i
	case int:
		return value.(int)
	default:
		return defaultValue
	}
}

func BoolToString(value bool, trueStr, falseStr string) string {
	if value {
		return trueStr
	} else {
		return falseStr
	}
}

func Quote(text string, quotationMark string) string {
	return quotationMark + text + quotationMark
}

func QuoteAndJoin(texts []string, quotationMark string, separator string) string {
	var elements []string
	for _, text := range texts {
		elements = append(elements, Quote(text, quotationMark))
	}
	return strings.Join(elements, separator)
}

func SQLValToInt(sqlVal *sqlparser.SQLVal, defaultValue int) int {
	if sqlVal == nil {
		return defaultValue
	}
	if value, err := strconv.Atoi(string(sqlVal.Val)); err != nil {
		return defaultValue
	} else {
		return value
	}
}
func SQLValToString(sqlVal *sqlparser.SQLVal, defaultValue string) string {
	if sqlVal == nil {
		return defaultValue
	}
	return string(sqlVal.Val)
}

func IfThenElseString(condition bool, trueValue, falseValue string) string {
	if condition {
		return trueValue
	} else {
		return falseValue
	}
}
func IfThenElseBool(condition bool, trueValue, falseValue bool) bool {
	if condition {
		return trueValue
	} else {
		return falseValue
	}
}
func IfThenElseFloat64(condition bool, trueValue, falseValue float64) float64 {
	if condition {
		return trueValue
	} else {
		return falseValue
	}
}

// ToLowerCamel converts snakeCase to camelCase.
// returns false if string conversion is insymmetric.
func ToLowerCamel(s string) (string, bool) {
	camel := strcase.ToLowerCamel(s)
	snake := strcase.ToSnake(camel)
	return camel, s == snake
}

// ToUpperCamel converts snakeCase to camelCase.
// returns false if string conversion is insymmetric.
func ToUpperCamel(s string) (string, bool) {
	camel := strcase.ToCamel(s)
	snake := strcase.ToSnake(camel)
	return camel, s == snake
}

// ToUpperCamel converts snakeCase to lower snakeCase.
// returns false if string conversion is insymmetric.
func ToLowerSnake(s string) (string, bool) {
	snake := strcase.ToSnake(strings.ToLower(s))
	return snake, s == snake
}

func NewTemplate(tplName, tplText string, funcMap template.FuncMap) (*template.Template, error) {
	return template.New(tplName).Funcs(funcMap).Parse(tplText)
}
