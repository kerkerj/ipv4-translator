package main

import (
	"testing"
)

func Test_IPEncode(t *testing.T) {
	t.Run("H", func(t *testing.T) {
		encoded := IPEncode("127.0.0.1")

		if encoded != 2130706433 {
			t.Fail()
		}
	})
}
