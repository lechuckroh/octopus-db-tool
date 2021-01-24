package util

func NewStringMemoizer(fn func(interface{}) string) func(interface{}) string {
	cache := make(map[interface{}]string)

	return func(key interface{}) string {
		if value, found := cache[key]; found {
			return value
		}
		value := fn(key)
		cache[key] = value
		return value
	}
}
