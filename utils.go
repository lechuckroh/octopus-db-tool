package main

import (
	"github.com/google/go-cmp/cmp"
	"github.com/iancoleman/strcase"
	"github.com/xwb1989/sqlparser"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

func IsStringType(typ string) bool {
	lowerType := strings.ToLower(typ)
	stringTypes := []string{"char", "string"}
	for _, stringType := range stringTypes {
		if strings.Contains(lowerType, stringType) {
			return true
		}
	}
	return false
}

func IsBooleanType(typ string) bool {
	lowerType := strings.ToLower(typ)
	booleanTypes := []string{"bool", "boolean"}
	for _, booleanType := range booleanTypes {
		if lowerType == booleanType {
			return true
		}
	}
	return false
}

func IsNumericType(typ string) bool {
	lowerType := strings.ToLower(typ)
	numericTypes := []string{"decimal", "float", "double", "long", "bigint", "int", "smallint", "number"}
	for _, numericType := range numericTypes {
		if lowerType == numericType {
			return true
		}
	}
	return false
}

func IsIntType(typ string) bool {
	lowerType := strings.ToLower(typ)
	intTypes := []string{"long", "bigint", "int", "smallint"}
	for _, numericType := range intTypes {
		if lowerType == numericType {
			return true
		}
	}
	return false
}

func IsDateType(typ string) bool {
	lowerType := strings.ToLower(typ)
	dateTypes := []string{"date", "datetime"}
	for _, dateType := range dateTypes {
		if lowerType == dateType {
			return true
		}
	}
	return false
}

// ParseType parses column type
// returns name, size, scale
func ParseType(str string) (string, uint16, uint16) {
	r := regexp.MustCompile(`(?m)([a-zA-Z]+)\(([0-9]+)[\s,]*([0-9]*)\)`)
	matches := r.FindStringSubmatch(str)
	matchLen := len(matches)

	var name string
	var size, scale int

	if matchLen >= 2 {
		name = matches[1]

		if matchLen >= 3 {
			size, _ = strconv.Atoi(matches[2])
		}
		if matchLen == 4 {
			scale, _ = strconv.Atoi(matches[3])
		}
	} else {
		name = str
	}

	return name, uint16(size), uint16(scale)
}

func GetFileFormat(fileFormat string, filename string) string {
	if fileFormat != "" {
		return fileFormat
	}

	ext := filepath.Ext(filename)
	switch strings.ToLower(ext) {
	case ".graphql":
		fallthrough
	case ".graphqls":
		return FormatGraphql
	case ".mdj":
		return FormatStaruml2
	case ".ojson":
		return FormatOctopus
	case ".plantuml":
		return FormatPlantuml
	case ".schema":
		return FormatSchemaConverter
	case ".xlsx":
		return FormatXlsx
	default:
		return ""
	}
}

func WriteLinesToFile(filename string, lines []string) error {
	if err := ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return err
	}
	log.Printf("[WRITE] %s", filename)
	return nil
}


func Quote(text string, quotationMark string) string {
	return quotationMark + text + quotationMark
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

type StringSet struct {
	valueMap map[string]bool
}

func NewStringSet(items ...string) *StringSet {
	set := &StringSet{
		valueMap: make(map[string]bool),
	}
	for _, item := range items {
		set.Add(item)
	}
	return set
}

func (s *StringSet) Add(value string) {
	s.valueMap[value] = true
}

func (s *StringSet) AddAll(values []string) {
	for _, value := range values {
		s.valueMap[value] = true
	}
}

func (s *StringSet) Contains(key string) bool {
	_, ok := s.valueMap[key]
	return ok
}

func (s *StringSet) ContainsAny(keys []string) bool {
	for _, key := range keys {
		if _, ok := s.valueMap[key]; ok {
			return true
		}
	}
	return false
}

func (s *StringSet) Remove(key string) {
	delete(s.valueMap, key)
}

func (s *StringSet) Clear() {
	s.valueMap = make(map[string]bool)
}

func (s *StringSet) Slice() []string {
	keys := make([]string, 0, len(s.valueMap))
	for key := range s.valueMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (s *StringSet) Join(separator string) string {
	return strings.Join(s.Slice(), separator)
}

func (s *StringSet) Equals(other *StringSet) bool {
	return cmp.Equal(s.valueMap, other.valueMap)
}

func (s *StringSet) Size() int {
	return len(s.valueMap)
}

func TernaryString(condition bool, trueValue, falseValue string) string {
	if condition {
		return trueValue
	} else {
		return falseValue
	}
}
func TernaryBool(condition bool, trueValue, falseValue bool) bool {
	if condition {
		return trueValue
	} else {
		return falseValue
	}
}
func TernaryFloat64(condition bool, trueValue, falseValue float64) float64 {
	if condition {
		return trueValue
	} else {
		return falseValue
	}
}

func NewBool(value bool) *bool {
	b := value
	return &b
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

func ToLowerSnake(s string) string {
	return strcase.ToSnake(strings.ToLower(s))
}
