package util

import (
	"github.com/google/go-cmp/cmp"
	"sort"
	"strings"
)

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

func (s *StringSet) ContainsAll(keys []string) bool {
	for _, key := range keys {
		if _, ok := s.valueMap[key]; !ok {
			return false
		}
	}
	return true
}

func (s *StringSet) Remove(key string) {
	delete(s.valueMap, key)
}

func (s *StringSet) RemoveAll(keys []string) {
	for _, key := range keys {
		delete(s.valueMap, key)
	}
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


