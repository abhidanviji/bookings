package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows does not have required fields when it does")
	}

}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	isValid := form.Has("a")
	if isValid {
		t.Error("form shows valid when it does not have the field")
	}

	postedData.Add("a", "a")

	form = New(postedData)
	isValid = form.Has("a")
	if !isValid {
		t.Error("form shows invalid when the field exists")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	isValid := form.MinLength("a", 1)
	if isValid {
		t.Error("form shows valid when the field does not exist")
	}

	postedData.Add("a", "a")

	form = New(postedData)

	isValid = form.MinLength("a", 1)
	if !isValid {
		t.Error("form shows invalid when the field has value with MinLength")
	}

	isValid = form.MinLength("a", 2)
	if isValid {
		t.Error("form shows valid when the field has value lesser than the MinLength")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("form shows valid when there is no email field")
	}

	postedData.Add("email", "a@a.com")
	postedData.Add("notemail", "a")
	form = New(postedData)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("form shows invalid when there is a valid email")
	}

	form.IsEmail("notemail")

	if form.Valid() {
		t.Error("form shows valid when the field is not email")
	}
}
