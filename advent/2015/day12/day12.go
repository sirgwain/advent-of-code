package day12

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {

	var s1, s2 int

	in := bytes.TrimSpace(input)
	s1 = sumNumbers(in)
	s2, err := sumNumbersJsonDecoder(in)
	if err != nil {
		return 0, 0, err
	}

	return s1, s2, nil
}

// sumNumbers iterates through in backwards summing
// up all numbers found in the string
func sumNumbers(in []byte) (sum int) {

	var cur int
	tens := 1
	// go through backwards so it's easier to
	// process numbers like -123
	for i := len(in) - 1; i >= 0; i-- {
		b := in[i]
		if b >= '0' && b <= '9' {
			cur += int(b-'0') * tens
			tens *= 10
		} else {
			// if we hit a negative on our way through
			// negate the current number (noop for 0)
			if b == '-' {
				cur *= -1
			}
			// add the number to our sum (noop for 0) and reset
			sum += cur
			cur = 0
			tens = 1
		}
	}

	return sum
}

// sum numbers excluding objects with a property of "red"
func sumNumbersJsonDecoder(in []byte) (int, error) {
	// create a new json decoder to walk
	dec := json.NewDecoder(bytes.NewReader(in))
	dec.UseNumber()

	type kind uint8
	const (
		object kind = iota
		array
	)

	type frame struct {
		kind   kind
		sum    int
		hasRed bool
	}
	stack := []frame{{kind: array}}
	var f frame

	for {
		token, err := dec.Token()
		if err != nil {
			// io.EOF when done
			if errors.Is(err, io.EOF) {
				return stack[0].sum, nil
			}
			return 0, fmt.Errorf("invalid json input %v", err)
		}
		switch v := token.(type) {
		case json.Delim:
			switch v {
			case '{':
				// pop object onto stack
				stack = append(stack, frame{kind: object})
			case '[':
				// pop array onto stack
				stack = append(stack, frame{kind: array})
			case '}':
				// pop object off stack
				f = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if !f.hasRed {
					stack[len(stack)-1].sum += f.sum
				}
			case ']':
				// pop array off stack
				f := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				// add the sum of the array frame
				stack[len(stack)-1].sum += f.sum
			}
		case string:
			top := &stack[len(stack)-1]
			if v == "red" && top.kind == object {
				top.hasRed = true
			}
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				return 0, fmt.Errorf("invalid json number: %v", err)
			}

			stack[len(stack)-1].sum += int(n)
		}
	}
}
