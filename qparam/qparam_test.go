package qparam

import "testing"

func TestEncode(t *testing.T) {
	type testStruct struct {
		Bool     bool   `qparam:"bool"`
		Str      string `qparam:"str"`
		EmptyStr string `qparam:"empty_str"`
		Int      int    `qparam:"int"`
		Float    float64
		Ignore   string `qparam:"-"`
	}

	vals, err := Marshal(testStruct{
		Bool:   true,
		Str:    "hello world",
		Int:    123,
		Float:  123.456,
		Ignore: "ignored",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if vals.Get("bool") != "true" {
		t.Errorf("expected bool to be true, got %s", vals.Get("bool"))
	}

	if _, ok := vals["empty_str"]; ok {
		t.Errorf("expected empty_str to be not be present in vals, but it was")
	}

	if vals.Get("str") != "hello world" {
		t.Errorf("expected str to be hello world, got %s", vals.Get("str"))
	}

	if vals.Get("int") != "123" {
		t.Errorf("expected int to be 123, got %s", vals.Get("int"))
	}

	if vals.Get("Float") != "" {
		t.Errorf("expected float to be empty, got %s", vals.Get("float"))
	}

	if vals.Get("Ignore") != "" {
		t.Errorf("expected ignore to be empty, got %s", vals.Get("ignore"))
	}
}
