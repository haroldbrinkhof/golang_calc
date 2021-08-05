package lexer

import (
	"strconv"
	"strings"
)

type TokenType byte
type OperatorType float64

type Chain struct {
	first  *Token
	input  string
	length uint32
}

func (c *Chain) Length() uint32 {
	return c.length
}

func (c *Chain) FirstToken() *Token {
	return c.first
}

type Token struct {
	value float64
	type_ TokenType
	next  *Token
}

func (t *Token) NextToken() *Token {
	return t.next
}
func (t *Token) Type() TokenType {
	return t.type_
}
func (t *Token) Value() float64 {
	return t.value
}

const operators string = "+-*/()"
const opFirst OperatorType = -1
const OpNone OperatorType = 0
const OpPlus OperatorType = 1
const OpMinus OperatorType = 2
const OpMultiply OperatorType = 3
const OpDivide OperatorType = 4
const OpParensLeft OperatorType = 5
const OpParensRight OperatorType = 6

func (o OperatorType) Priority() int8 {
	switch o {
	case OpPlus:
		fallthrough
	case OpMinus:
		return 0
	case OpDivide:
		fallthrough
	case OpMultiply:
		return 1
	case OpParensLeft:
		fallthrough
	case OpParensRight:
		return 3
	default:
		return -1
	}
}

const Numeral TokenType = 1
const Operator TokenType = 2

func Process(input string) (*Chain, error) {
	var chain Chain
	sanitized := sanitizeString(input)

	chain = Chain{nil, input, 0}
	if len(sanitized) == 0 {
		return &chain, nil
	}

	var t *Token
	for prevOperator := opFirst; len(sanitized) > 0; {
		if chain.first == nil {
			t = &Token{0.0, 0, nil}
			chain.first = t
		} else {
			t.next = &Token{0.0, 0, nil}
			t = t.next
		}
		value, operatorFound, hasValue, remainder, err := getNextToken(sanitized, prevOperator)
		prevOperator = operatorFound
		if err != nil {
			return nil, err
		}
		sanitized = remainder
		if hasValue {
			t.value = value
			t.type_ = Numeral
			t.next = nil
		} else {
			t.value = float64(prevOperator)
			t.type_ = Operator
			t.next = nil
		}
		chain.length++

	}
	return &chain, nil
}

func getNextToken(remainder string, prevOperator OperatorType) (float64, OperatorType, bool, string, error) {
	p := strings.IndexAny(remainder, operators)
	if p >= 0 {
		op := determineOperator([]rune(remainder)[p])
		isSign := p == 0 && len(remainder) > 1 && (op == OpMinus || op == OpPlus) && prevOperator != OpNone && prevOperator != OpParensRight

		if isSign {
			// we have a + or - sign indicator, not operator, get value following
			remainder = remainder[1:]
			// get the position of the next operator
			p = strings.IndexAny(remainder, operators)
			if p == -1 {
				// if no more operator is found make sure slicing still works correctly
				p = len(remainder)
			}
			// get the value up to the next operator
			value, hasValue := determineValue(remainder[:p])
			if op == OpMinus && hasValue {
				// if our first operator is a minus negate the value
				value = -value
			}
			return value, OpNone, true, remainder[p:], nil
		} else if p > 0 {
			// we have a value in everything preceeding p
			value, _ := determineValue(remainder[:p])
			return value, OpNone, true, remainder[p:], nil

		} else {
			// we found an operator
			return 0.0, op, false, remainder[p+1:], nil
		}
	} else {
		// no more operator found, get remaining value
		value, hasValue := determineValue(remainder)
		return value, OpNone, hasValue, "", nil
	}
}

func determineValue(remainder string) (float64, bool) {
	value, hasValue := 0.0, false
	if len(remainder) > 0 {
		var err error
		value, err = strconv.ParseFloat(remainder, 64)
		if err == nil {
			hasValue = true
		}
	}
	return value, hasValue
}

func determineOperator(op rune) OperatorType {
	switch {
	case op == '-':
		return OpMinus
	case op == '+':
		return OpPlus
	case op == '*':
		return OpMultiply
	case op == '/':
		return OpDivide
	case op == '(':
		return OpParensLeft
	case op == ')':
		return OpParensRight
	default:
		return OpNone
	}
}

func sanitizeString(input string) string {
	killWhiteSpace := func(r rune) rune {
		switch {
		case r == ' ':
			fallthrough
		case r == '\t':
			fallthrough
		case r == '\n':
			return -1
		default:
			return r
		}
	}

	return strings.Map(killWhiteSpace, input)
}
