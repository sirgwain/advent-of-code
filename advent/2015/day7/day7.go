package day7

import (
	"bytes"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"unicode"
)

type Day struct {
}

type op string

const (
	ASSIGN op = "ASSIGN"
	AND    op = "AND"
	OR     op = "OR"
	NOT    op = "NOT"
	LSHIFT op = "LSHIFT"
	RSHIFT op = "RSHIFT"
)

type instruction struct {
	operand1 string
	operator op
	operand2 string
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	// build a list of instructions from the input
	instructions := make(map[string]instruction)
	for line := range strings.SplitSeq(string(bytes.TrimSpace(input)), "\n") {
		ins, out, err := parseLine(line)
		if err != nil {
			return 0, 0, err
		}
		instructions[out] = ins
	}

	// evaluate to determine the output on wire "a"
	s1 = evaluate("a", instructions, make(map[string]int))
	
	// for s2, take the signal on wire a, override wire b to that signal, reset
	instructions["b"] = instruction{
		operand1: strconv.Itoa(s1),
		operator: ASSIGN,
	}
	s2 = evaluate("a", instructions, make(map[string]int))

	return s1, s2, nil
}

// evaluate recursively evaluates instructions to find the output on a wire
// results are cached
func evaluate(wire string, instructions map[string]instruction, cache map[string]int) int {
	if v, ok := cache[wire]; ok {
		// already evaluated this wire
		return v
	}
	ins := instructions[wire]
	slog.Debug("evaluting instruction", "op", ins.operator, "o1", ins.operand1, "o2", ins.operand2, "out", wire)

	// find any literal operands of this instruction
	var operand1 int = -1
	if unicode.IsDigit(rune(ins.operand1[0])) {
		var err error
		operand1, err = strconv.Atoi(ins.operand1)
		if err != nil {
			panic(fmt.Errorf("failed to parse literal from operand: %s %v", ins.operand1, err))
		}
	}
	var operand2 int = -1
	if ins.operand2 != "" && unicode.IsDigit(rune(ins.operand2[0])) {
		var err error
		operand2, err = strconv.Atoi(ins.operand2)
		if err != nil {
			panic(fmt.Errorf("failed to parse literal from operand: %s %v", ins.operand2, err))
		}
	}

	if operand1 == -1 {
		operand1 = evaluate(ins.operand1, instructions, cache)
	}
	if operand2 == -1 && (ins.operator == AND || ins.operator == OR || ins.operator == LSHIFT || ins.operator == RSHIFT) {
		operand2 = evaluate(ins.operand2, instructions, cache)
	}

	// evaluate this instruction based on the operator
	var out int
	switch ins.operator {
	case ASSIGN:
		return operand1
	case NOT:
		out = ^operand1
	case AND:
		out = operand1 & operand2
	case OR:
		out = operand1 | operand2
	case LSHIFT:
		out = operand1 << operand2
	case RSHIFT:
		out = operand1 >> operand2
	default:
		panic("unknonwn operation")
	}

	cache[wire] = out
	return out
}

// parseLine parses a single instruction line returning the instruction and the output wire
func parseLine(line string) (instruction, string, error) {
	// each input is some form of
	// identifier -> wire
	// 456 -> y
	// x AND y -> z
	// NOT y -> i
	// 1 AND x -> y
	// x LSHIFT 2 -> y
	// x LSHIFT y -> z

	splits := strings.Split(line, " -> ")
	fields := strings.Fields(splits[0])
	out := splits[1]

	switch len(fields) {
	case 1:
		// 456 -> y
		// lx -> a
		return instruction{
			operand1: fields[0],
			operator: ASSIGN,
		}, out, nil
	case 2:
		// NOT y -> i
		if fields[0] != "NOT" {
			return instruction{}, "", fmt.Errorf("unknown unary operator from line: %s", line)
		}
		return instruction{
			operator: NOT,
			operand1: fields[1],
		}, out, nil
	case 3:
		// x AND y -> z
		return instruction{
			operator: op(fields[1]),
			operand1: fields[0],
			operand2: fields[2],
		}, out, nil
	}
	return instruction{}, "", fmt.Errorf("unable to build instruction from line %s", line)
}
