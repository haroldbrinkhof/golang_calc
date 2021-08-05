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
