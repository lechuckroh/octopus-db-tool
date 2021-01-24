package util

import (
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

func TestNewStringMemoizer(t *testing.T) {
	value := 0

	fn := func(key interface{}) string {
		value += 1
		return strconv.Itoa(value)
	}
	cache := NewStringMemoizer(fn)

	Convey("", t, func() {
		So(cache("1"), ShouldEqual, "1")
		So(cache("2"), ShouldEqual, "2")
		So(cache("3"), ShouldEqual, "3")
		So(cache("1"), ShouldEqual, "1")
		So(fn("1"), ShouldEqual, "4")
	})
}
