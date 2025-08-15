package json_writer

import (
	"testing"
)

func TestAreDifferent(t *testing.T) {
	if End == Obj {
		t.Error("Obj should not be equal to Err")
	}
	if Obj == Arr {
		t.Error("Obj should not be equal to Err")
	}
	if Arr == End {
		t.Error("Obj should not be equal to Err")
	}
}
