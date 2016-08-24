package shared

import (
	"testing"
)

func Test_SnakeCase(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{in: "testing-string", out: "testing_string"},
		{in: "TestingString", out: "testing_string"},
		{in: "Testing_String", out: "testing__string"},
		{in: "", out: ""},
	}

	for _, test := range cases {
		if out := SnakeCase(test.in); out != test.out {
			t.Errorf("expected %s but got %s", test.out, out)
		}
	}
}

func Test_ParameterizeString(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{in: "testing-string", out: "testing_string"},
		{in: "TestingString", out: "testingstring"},
		{in: "Testing-String", out: "testing_string"},
		{in: "", out: ""},
	}

	for _, test := range cases {
		if out := ParameterizeString(test.in); out != test.out {
			t.Errorf("expected %s but got %s", test.out, out)
		}
	}
}
