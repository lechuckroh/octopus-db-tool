package main

import (
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