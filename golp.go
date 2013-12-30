package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Evaluate an expression
func eval(e interface{}) interface{} {
	switch e := e.(type) {
	case bool, uint64, float64:
		return e
	case string:
		return e
	case []string:
		switch e[0] {
		case "quote": // (quote exp)
			return e[1:]
		case "if": // (if test then else?)
			test := e[1]
			then := e[2]
			else_ := e[3]
			if eval(test).(bool) {
				return eval(then)
			} else {
				return eval(else_)
			}
			// case "def": // (def var exp)
			// case "fn":
			// case "set!": // (set! var exp)
		}
	}
	return e
}

// Read a Scheme expression from a string.
func read(s string) ([]interface{}, error) {
	return (read_from(tokenize(s)))
}

// Converts a string into an array of tokens.
func tokenize(s string) []string {
	return regexp.MustCompile(`\s+`).Split(
		strings.Replace(strings.Replace(s, "(", " ( ", -1), ")", " ) ", -1), -1)
}

// Read an expression from a sequence of tokens.
func read_from(tokens []string) ([]interface{}, error) {
	var token string
	if len(tokens) == 0 {
		return nil, errors.New("unexpected EOF while reading")
	}
	token, tokens = tokens[len(tokens)-1], tokens[:len(tokens)-1]
	switch token {
	case "(":
		l := make([]interface{}, 0)
		for tokens[0] != ")" {
			token, _ := read_from(tokens) // TODO handle err
			l = append(l, token)
		}
		tokens = tokens[:len(tokens)-1]
		return l, nil
	case ")":
		return nil, errors.New("unexpected )")
	}
	return []interface{}{atom(token)}, nil
}

// Bools, ints, and floats are converted; every other token is a symbol.
func atom(s string) interface{} {
	if b, err := strconv.ParseBool(s); (s == "true" || s == "false") && err != nil {
		return b
	}
	if i, err := strconv.ParseUint(s, 0, 64); err != nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err != nil {
		return f
	}
	return s
}

func prompt() {
	fmt.Print(">>> ")
}

func main() {
	//env := &Env{}

	scanner := bufio.NewScanner(os.Stdin)
	for prompt(); scanner.Scan(); prompt() {
		in := scanner.Text()
		fmt.Println(in)
		parsed, _ := read(in) // TODO handle err
		fmt.Println(eval(parsed))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
