package util

import (
	"regexp"
	"strconv"
	"strings"
)

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
