package util

import (
	"github.com/iancoleman/strcase"
	"testing"
)

func TestStringSet_Equals(t *testing.T) {
	s1 := NewStringSet()
	s2 := NewStringSet()

	s1.AddAll([]string {"foo", "bar"})
	s2.AddAll([]string {"bar", "foo"})

	if !s1.Equals(s2) {
		t.Error("TestStringSet_Equals #1 failed")
	}

	s2.Add("foobar")
	if s1.Equals(s2) {
		t.Error("TestStringSet_Equals #2 failed")
	}
}

func TestToLowerCamel(t *testing.T) {
	trueCase := []string{"id", "user_name"}
	for _, s := range trueCase {
		camel, ok := ToLowerCamel(s)
		if !ok {
			t.Errorf("ToLowerCamel failed: %s -> %s -> %s", s, camel, strcase.ToSnake(camel))
		}
	}

	falseCase := []string{"a1_code", "a1"}
	for _, s := range falseCase {
		camel, ok := ToLowerCamel(s)
		if ok {
			t.Errorf("ToLowerCamel failed: %s -> %s -> %s", s, camel, strcase.ToSnake(camel))
		}
	}
}

func TestToLowerSnake(t *testing.T) {
	trueCase := []string{"id", "user_name"}
	for _, s := range trueCase {
		snake, ok := ToLowerSnake(s)
		if !ok {
			t.Errorf("ToLowerSnake failed: %s -> %s", s, snake)
		}
	}

	falseCase := []string{"a1_code", "value_level1"}
	for _, s := range falseCase {
		snake, ok := ToLowerSnake(s)
		if ok {
			t.Errorf("ToLowerSnake failed: %s -> %s", s, snake)
		}
	}
}
