package parser

import (
	"testing"
)

func TestParse_emptyInput_chainLengthOf0(t *testing.T) {
	if expected := sanitizeString(""); expected != "" {
		t.Errorf("empty input should have returned nil")
	}
}

func TestParse_whitespaceGetsRemoved(t *testing.T) {
	if expected := sanitizeString(" 1\t +3 4\n"); expected != "1+34" {
		t.Errorf("all whitespace should have been removed but we ended up with '%s'", expected)
	}
}

func TestParse(t *testing.T) {
	var expected *Chain
	expected, _ = Parse("+1 + -34")
	if expected.length != 3 {
		t.Errorf("3 tokens should have been found but we found only %d", expected.length)
	} else {
		token := expected.first
		if token.type_ != Numeral || token.value != 1.0 {
			t.Errorf("expected numeral with value 1 but got %f", token.value)

		}
		if token = token.next; token.type_ != Operator || token.value != 1.0 {
			t.Errorf("expected plus operator %f but got %f", OpPlus, token.value)

		}
		if token = token.next; token.type_ != Numeral || token.value != -34.0 || token.next != nil {
			t.Errorf("expected numeral with value -34 but got %f", token.value)

		}
	}

}
