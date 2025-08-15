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

func TestEmptyObject(t *testing.T) {
	j := New()
	err := j.Write(Obj)
	if err != nil {
		t.Fatalf("Error writing start of object: %v", err)
	}
	err = j.Write(End)
	if err != nil {
		t.Fatalf("Error writing end of object: %v", err)
	}
	if j.String() != "{}" {
		t.Errorf("Expected '{}', got '%s'", j.String())
	}
}
