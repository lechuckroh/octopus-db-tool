package main

import (
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

func GetDefaultString(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
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
	numericTypes := []string{"float", "double", "long", "bigint", "int", "smallint", "number"}
	for _, numericType := range numericTypes {
		if lowerType == numericType {
			return true
		}
	}
	return false
}

func ParseType(str string) (string, uint16) {
	r := regexp.MustCompile(`(?m)([a-zA-Z]+)\(([0-9]+)\)`)
	matches := r.FindStringSubmatch(str)
	if len(matches) == 3 {
		size, _ := strconv.Atoi(matches[2])
		return matches[1], uint16(size)
	}
	return str, 0
}

func GetFileFormat(fileFormat string, filename string) string {
	if fileFormat != "" {
		return fileFormat
	}

	ext := filepath.Ext(filename)
	switch strings.ToLower(ext) {
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


type StringSet struct {
	valueMap map[string]bool
}

func NewStringSet() *StringSet {
	return &StringSet{
		valueMap: make(map[string]bool),
	}
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

func (s *StringSet) Slice() []string {
	keys := make([]string, 0, len(s.valueMap))
	for key := range s.valueMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
