package main

import (
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
