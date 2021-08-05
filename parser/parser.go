package parser

import (
	"errors"

	"catsandcoding.be/calc/lexer"
)

type Term struct {
	leftHand  *float64
	rightHand *float64
	operator  lexer.OperatorType
	leftTerm  *Term
	rightTerm *Term
	priority  uint16
}

type Chain struct {
	first  *Term
	length uint32
}

func Calculate(input string) (float64, error) {
	lexerChain, _ := lexer.Process(input)
	parserChain, _ := transform(lexerChain)
	return parse(parserChain)
}

func parse(chain *Chain) (float64, error) {
	return 0, nil
}

func transform(chain *lexer.Chain) (*Chain, error) {
	outcome := new(Chain)

	if chain.Length() > 0 {
		var basePriority uint16 = 0
		var t *lexer.Token = chain.FirstToken()
		var term *Term = new(Term)
		outcome.first = term
		term.operator = lexer.OpNone

		for t != nil {
			if t.Type() == lexer.Operator {
				// we have an operator determine priority
				switch {
				case t.Value() == float64(lexer.OpParensLeft):
					basePriority += 3
				case t.Value() == float64(lexer.OpParensRight):
					basePriority -= 3
					if basePriority > 0 {
						return nil, errors.New("Parse error: mismatching parentheses")
					}
				default:
					term.operator = lexer.OperatorType(t.Value())
					term.priority = uint16(basePriority) + uint16(term.operator.Priority())
					// if no value has been set yet, steal the right one from the previous term
					if term.leftHand == nil {
						if term.leftTerm == nil {
							return nil, errors.New("Parse error: numerals must preceed an operator")
						}
						term.leftHand = term.leftTerm.rightHand
						term.leftTerm.rightHand = nil
					}
				}
			} else {
				// we have a value, store it left or right
				if term.leftHand == nil {
					//fmt.Printf("TERM LEFTHAND\n")
					value := t.Value()
					term.leftHand = &value
				} else {
					//fmt.Printf("TERM RIGHTHAND\n")
					// we must have an operator already at this point
					if term.operator == lexer.OpNone {
						return nil, errors.New("Parse error: 2 or more successive numerals without operator")
					}
					//fmt.Printf("TERM RIGHTHAND not OpNone\n")
					value := t.Value()
					term.rightHand = &value

					// create the next term
					if t.NextToken() != nil {
						nextTerm := new(Term)
						term.rightTerm = nextTerm
						nextTerm.leftTerm = term
						nextTerm.rightTerm = nil
						nextTerm.operator = lexer.OpNone
						term = nextTerm
						outcome.length++
					}
				}
			}

			t = t.NextToken()

		}
		// if we end our input with parentheses we get an empty term, remove it
		if term.leftHand == nil && term.rightHand == nil && term.operator == lexer.OpNone {
			term.leftTerm.rightTerm = nil
		}

	} else {
		outcome.first = nil
		outcome.length = 0
	}

	return outcome, nil
}
