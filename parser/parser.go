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

func (t *Term) calculate() (float64, error) {
	if t.leftHand == nil || t.rightHand == nil {
		return 0, errors.New("Parse error can not calculate total; both left hand and right hand need to be filled in.")
	}

	switch t.operator {
	case lexer.OpPlus:
		total := *(t.leftHand) + *(t.rightHand)
		return total, nil
	case lexer.OpMinus:
		total := *(t.leftHand) - *(t.rightHand)
		return total, nil
	case lexer.OpDivide:
		if *(t.rightHand) == 0 {
			return 0, errors.New("Parse error: attempted division by zero")
		}
		total := *(t.leftHand) / *(t.rightHand)
		return total, nil
	case lexer.OpMultiply:
		total := *(t.leftHand) * *(t.rightHand)
		return total, nil
	default:
		return 0, errors.New("Parse error can not calculate total of this term; operator.")
	}
}

type Chain struct {
	first  *Term
	length uint32
}

func Calculate(input string) (float64, error) {
	lexerChain, err := lexer.Process(input)
	if err != nil {
		return 0, err
	}
	parserChain, err := transform(lexerChain)
	if err != nil {
		return 0, err
	}
	return parse(parserChain)
}

func parse(chain *Chain) (float64, error) {
	outcome := 0.0
	term := chain.first

	for term != nil {
		switch {
		case term.rightTerm != nil && term.priority < term.rightTerm.priority:
			term = term.rightTerm
			continue
		default:
			// make sure left and right hand are filled in
			// by stealing them from neighbouring terms if necessary
			err := stealValues(term)
			if err != nil {
				return 0, err
			}

			// we have a complete term, calculate it
			total, err := term.calculate()
			if err != nil {
				return 0, err
			}

			// with the total calculated we now assign it to
			// the nearest appropriate term,i.e the one with highest priority
			// or the left side
			if term.rightTerm != nil && term.rightTerm.priority == term.priority {
				// right side has priority assign here
				if term.rightTerm.leftHand == nil {
					term.rightTerm.leftHand = &total
				}

				// after we assigned the value we remove the current
				// term from the chain by adjusting the pointers of its
				// neighbours
				if term.leftTerm != nil {
					term.leftTerm.rightTerm = term.rightTerm
				}
				term.rightTerm.leftTerm = term.leftTerm

				// adjustment might lead to no term available on the left,
				// if so then this is our chain's new start
				if term.leftTerm == nil {
					chain.first = term.rightTerm
				}

				// start over from the beginning of the chain
				term = chain.first
				continue
			} else if term.leftTerm != nil {
				// we have a left term, assign here
				if term.leftTerm.rightHand == nil {
					term.leftTerm.rightHand = &total
				}

				// after we assigned the value we remove the current
				// term from the chain by adjusting the pointers of its
				// neighbours
				if term.rightTerm != nil {
					term.rightTerm.leftTerm = term.leftTerm
				}
				term.leftTerm.rightTerm = term.rightTerm

				// adjustment might lead to no term available on the left,
				// if so then this is our chain's new start
				if term.leftTerm == nil {
					chain.first = term.rightTerm
				}

				// start over from the beginning of the chain
				term = chain.first
				continue

			} else if term.leftTerm == nil && term.rightTerm == nil {
				// no neighbours left, everything calculated
				term = nil
				outcome = total
			}

		}
	}

	return outcome, nil
}
func stealValues(t *Term) error {
	if t.leftHand == nil {
		if t.leftTerm == nil || t.leftTerm.rightHand == nil {
			return errors.New("Parse error: no value to steal for left hand")
		}
		t.leftHand = t.leftTerm.rightHand
		t.leftTerm.rightHand = nil
	}
	if t.rightHand == nil {
		if t.rightTerm == nil || t.rightTerm.leftHand == nil {
			return errors.New("Parse error: no value to steal for right hand")
		}
		t.rightHand = t.rightTerm.leftHand
		t.rightTerm.leftHand = nil
	}
	return nil
}

func transform(chain *lexer.Chain) (*Chain, error) {
	outcome := new(Chain)

	if chain.Length() > 0 {
		var basePriority int32 = 0
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
					if basePriority < 0 {
						return nil, errors.New("Parse error: mismatching parentheses, too many )")
					}
					// if we are other than ( or  ) assign operator as well
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
				if term.leftHand == nil && term.operator == lexer.OpNone {
					value := t.Value()
					term.leftHand = &value
				} else {
					// we must have an operator already at this point
					if term.operator == lexer.OpNone {
						return nil, errors.New("Parse error: 2 or more successive numerals without operator")
					}
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
		if basePriority > 0 {
			return nil, errors.New("Parse error: mismatching parentheses, too many (")
		}

	} else {
		outcome.first = nil
		outcome.length = 0
	}

	return outcome, nil
}
