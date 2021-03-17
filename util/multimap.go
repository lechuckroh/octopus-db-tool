package util

type MultiMap interface {
	Clear()
	ContainsKey(key interface{}) bool
	Empty() bool
	Get(key interface{}) (value []interface{}, found bool)
	Keys() []interface{}
	Put(key interface{}, value interface{})
	RemoveKey(key interface{})
	RemoveValue(key interface{}, value interface{})
}

type SMultiMap struct {
	store map[interface{}][]interface{}
}

func NewMultiMap() MultiMap {
	return &SMultiMap{store: make(map[interface{}][]interface{})}
}

func (m *SMultiMap) Clear() {
	m.store = make(map[interface{}][]interface{})
}

func (m *SMultiMap) ContainsKey(key interface{}) bool {
	_, found := m.store[key]
	return found
}

func (m *SMultiMap) Empty() bool {
	return len(m.store) == 0
}

func (m *SMultiMap) Get(key interface{}) ([]interface{}, bool) {
	values, found := m.store[key]
	return values, found
}

func (m *SMultiMap) Keys() []interface{} {
	keys := make([]interface{}, len(m.store))
	idx := 0
	for key, _ := range m.store {
		keys[idx] = key
		idx++
	}
	return keys
}

func (m *SMultiMap) Put(key interface{}, value interface{}) {
	m.store[key] = append(m.store[key], value)
}

func (m *SMultiMap) RemoveKey(key interface{}) {
	delete(m.store, key)
}

func (m *SMultiMap) RemoveValue(key interface{}, value interface{}) {
	values, found := m.store[key]
	if found {
		for i, v := range values {
			if v == value {
				m.store[key] = append(values[:i], values[i+1:]...)
			}
		}
	}
	if len(m.store[key]) == 0 {
		delete(m.store, key)
	}
}
