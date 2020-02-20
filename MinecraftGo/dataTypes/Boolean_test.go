package dataTypes

import (
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name  string
		value bool
	}{
		{name: "True", value: true},
		{name: "False", value: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			write := WriteBoolean(test.value)
			read, _ := ReadBoolean(write)
			if read != test.value {
				t.Fail()
			}
		})
	}
}
