package forms

import "testing"

func TestErrors_Add(t *testing.T) {
	var e errors = make(map[string][]string)
	e.Add("a", "Test error message")
	if len(e) != 1 {
		t.Error("there should be just one error message")
	}
}

func TestErrors_Get(t *testing.T) {
	var e errors = make(map[string][]string)
	expectedMsg := "Test error message for a"
	e.Add("a", expectedMsg)
	e.Add("b", "Test error message for b")
	msg := e.Get("a")
	if msg != expectedMsg {
		t.Error("the error message received does not match the expected")
	}

	msg = e.Get("c")
	if msg != "" {
		t.Error("received one error message when there shouldn't be one")
	}
}
