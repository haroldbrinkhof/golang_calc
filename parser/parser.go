package parser

import (
	"strconv"
	"strings"
)

type TokenType byte

type Chain struct {
	first  *Token
	input  string
	length int
}

type Token struct {
	value float64
	type_ TokenType
	next  *Token
}

const operators string = "+-*/()"
const opFirst float64 = -1
const opNone float64 = 0
const OpPlus float64 = 1
const OpMinus float64 = 2
const OpMultiply float64 = 3
const OpDivide float64 = 4
const OpParensLeft float64 = 5
const OpParensRight float64 = 6

const Numeral TokenType = 1
const Operator TokenType = 2

func Parse(input string) (*Chain, error) {
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
			t.value = prevOperator
			t.type_ = Operator
			t.next = nil
		}
		chain.length++

	}
	return &chain, nil
}

func getNextToken(remainder string, prevOperator float64) (float64, float64, bool, string, error) {
	p := strings.IndexAny(remainder, operators)
	if p >= 0 {
		op := determineOperator([]rune(remainder)[p])
		isSign := p == 0 && len(remainder) > 1 && (op == OpMinus || op == OpPlus) && prevOperator != opNone

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
			return value, opNone, true, remainder[p:], nil
		} else if p > 0 {
			// we have a value in everything preceeding p
			value, _ := determineValue(remainder[:p])
			return value, opNone, true, remainder[p:], nil

		} else {
			// we found an operator
			return 0.0, op, false, remainder[p+1:], nil
		}
	} else {
		// no more operator found, get remaining value
		value, hasValue := determineValue(remainder)
		return value, opNone, hasValue, "", nil
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

func determineOperator(op rune) float64 {
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
		return opNone
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
