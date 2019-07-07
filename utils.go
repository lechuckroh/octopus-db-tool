package main

import (
	"strconv"
)

func toInt(value interface{}, defaultValue int) int {
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
