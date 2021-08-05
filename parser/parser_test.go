package parser

import (
	"testing"

	"catsandcoding.be/calc/lexer"
)

func TestTransform(t *testing.T) {
	chain, _ := lexer.Process("3 + (5 * 3)")
	actual, err := transform(chain)
	if err != nil {
		t.Fatal(err)
	}
	if actual.length != 2 {
		t.Errorf("unexpected length, expected 2 but got %d\n", actual.length)
	}
	if term := actual.first; !(term.leftHand != nil && *(term.leftHand) == 3 && term.operator == lexer.OpPlus && term.rightHand == nil && term.leftTerm == nil && term.rightTerm != nil && term.priority == 0) {
		t.Fatal("incorrect first Term")

	}
	if term := actual.first.rightTerm; !(term.leftTerm == actual.first && term.rightTerm == nil && term.leftHand != nil && term.rightHand != nil && *(term.leftHand) == 5 && *(term.rightHand) == 3 && term.operator == lexer.OpMultiply && term.priority == 4) {
		t.Fatalf("incorrect second Term\nleftTerm(%p) == first(%p)\nrightTerm == nil(%t)\nleftHand != nil(%t)\nrightHand != nil(%t)\n", term.leftTerm, actual.first, term.rightTerm == nil, term.leftHand != nil, term.rightHand != nil)
	}
}

func TestParse(t *testing.T) {
	term1 := new(Term)
	term2 := new(Term)
	leftHand1 := 3.0

	term1.leftHand = &leftHand1
	term1.rightTerm = term2
	term1.operator = lexer.OpPlus
	term1.priority = 0

	leftHand2 := 5.0
	rightHand2 := 3.0
	term2.leftHand = &leftHand2
	term2.rightHand = &rightHand2
	term2.leftTerm = term1
	term2.operator = lexer.OpMultiply
	term2.priority = 4

	chain := new(Chain)
	chain.first = term1

	if actual, _ := parse(chain); actual != 18 {
		t.Errorf("expected 18 but got %f\n", actual)
	}

}

func TestCalculateMinus(t *testing.T) {

	if actual, _ := Calculate("1 - 2"); actual != -1 {
		t.Errorf("expected -1 but got %f\n", actual)
	}
}

func TestCalculatePlusNegative(t *testing.T) {

	if actual, _ := Calculate("1 + - 2"); actual != -1 {
		t.Errorf("expected -1 but got %f\n", actual)
	}
}

func TestCalculateMultiply(t *testing.T) {

	if actual,  := Calculate("3 * 2"); actual != 6 {
		t.Errorf("expected 6 but got %f\n", actual)
	}
}
func TestCalculateDivide(t *testing.T) {

	if actual, _ := Calculate("9 / 3"); actual != 3 {
		t.Errorf("expected 3 but got %f\n", actual)
	}
}
func TestCalculateDivideByZero(t *testing.T) {
	_, error := Calculate(" 9 / 0")
	if error == nil {
		t.Error("expected error")
	}
	if error.Error() != "Parse error: attempted division by zero" {
		t.Errorf("didn't expect %s\n", error.Error())
	}
}
func TestCalculate(t *testing.T) {
	if actual, _ := Calculate("1 - 2 * ((3 * 9) / 2) + 15"); actual != -11 {
		t.Errorf("expected -11 but got %f\n", actual)
	}
}

func TestParseMissingParentThesesParensRight(t *testing.T) {
	_, error := Calculate("1 + (3 * (1 + 7)")
	if error == nil {
		t.Fatal("expected error")
	}
	if error.Error() != "Parse error: mismatching parentheses, too many (" {
		t.Fatalf("expected different error but got %s\n", error)
	}

}
func TestParseMissingParentThesesParensLeft(t *testing.T) {
	_, error := Calculate("1 + (3 * 1 + 7))")
	if error == nil {
		t.Fatal("expected error")
	}
	if error.Error() != "Parse error: mismatching parentheses, too many )" {
		t.Fatalf("expected different error but got %s\n", error)
	}
}
