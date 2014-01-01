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

type Any interface{}

type Env struct {
	m     map[string]Any
	outer *Env
}

func (e *Env) find(k Any) (*Env, bool) {
	if _, ok := e.m[k.(string)]; ok {
		return e, true
	} else if e.outer != nil {
		return e.outer.find(k)
	}
	return nil, false
}

func (e *Env) get(k Any) Any {
	return e.m[k.(string)]
}

func (e *Env) merge(m map[string]Any) {
	for k, v := range m {
		e.m[k] = v
	}
}

func (e *Env) set(k Any, v Any) {
	e.m[k.(string)] = v
}

func NewEnv(keys []Any, vals []Any, outer *Env) *Env {
	zipped := map[string]Any{}
	vlen := len(vals)
	for i, key := range keys {
		if i < vlen {
			zipped[key.(string)] = vals[i]
		}
	}
	return &Env{zipped, outer}
}

// Evaluate an expression
func eval(e Any, env *Env) Any {
	switch e := e.(type) {
	case string:
		if environ, ok := env.find(e); ok {
			return environ.get(e)
		}
		return nil
	case bool, int64, float64:
		return e
	case []Any:
		if len(e) > 0 {
			switch e[0] {
			case "quote": // (quote exp)
				return e[1]
			case "if": // (if test then else?)
				if eval(e[1], env).(bool) {
					return eval(e[2], env)
				} else {
					return eval(e[3], env)
				}
			case "set!": // (set! var exp)
				if environ, ok := env.find(e[1]); ok {
					environ.set(e[1], eval(e[2], env))
				}
			case "define", "def": // (define var exp)
				env.set(e[1], eval(e[2], env))
			case "lambda", "fn": // (lambda (var*) exp)
				vars := e[1].([]Any)
				exp := e[2]
				return func(args ...Any) Any {
					return eval(exp, NewEnv(vars, args, env))
				}
			case "begin": // (begin exp*)
				var val Any
				for _, exp := range e[1:] {
					val = eval(exp, env)
				}
				return val
			default: // (proc exp*)
				exprs := make([]Any, len(e))
				for i, exp := range e {
					exprs[i] = eval(exp, env)
				}
				fn := exprs[0].(func(...Any) Any)
				return fn(exprs[1:]...)
			}
		}
	}
	return e
}

// Read a Scheme expression from a string.
func read(s string) (Any, error) {
	parsed, _, err := read_from(tokenize(s))
	return parsed, err
}

// Converts a string into an array of tokens.
func tokenize(s string) []string {
	return regexp.MustCompile(`\s+`).Split(strings.TrimSpace(
		strings.Replace(strings.Replace(s, "(", " ( ", -1), ")", " ) ", -1)), -1)
}

// Read an expression from a sequence of tokens.
func read_from(tokens []string) (Any, []string, error) {
	if len(tokens) == 0 {
		return nil, nil, errors.New("unexpected EOF while reading")
	}
	token := tokens[0]
	tokens = tokens[1:]
	switch token {
	case "(":
		L := []Any{}
		for len(tokens) > 0 && tokens[0] != ")" {
			token, remaining, _ := read_from(tokens) // TODO handle err
			L = append(L, token)
			tokens = remaining
		}
		if len(tokens) > 0 {
			tokens = tokens[1:] // pop off ')'
		}
		return L, tokens, nil
	case ")":
		return token, tokens, errors.New("unexpected )")
	}
	return atom(token), tokens, nil
}

// Bools, ints, and floats are converted; every other token is a symbol.
func atom(s string) Any {
	if b, err := strconv.ParseBool(s); (s == "true" || s == "false") && err == nil {
		return b
	}
	if i, err := strconv.ParseInt(s, 0, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

func prompt() {
	fmt.Print(">>> ")
}

func main() {
	env := NewEnv(nil, nil, nil)
	env.merge(map[string]Any{
		"+": func(xs ...Any) Any {
			sum := int64(0)
			for _, x := range xs {
				sum += x.(int64)
			}
			return sum
		},
		"-": func(xs ...Any) Any {
			switch len(xs) {
			case 0:
				return 0
			case 1:
				return -1 * xs[0].(int64)
			case 2:
				return xs[0].(int64) - xs[1].(int64)
			default:
				sum := xs[0].(int64) - xs[1].(int64)
				for _, x := range xs[2:] {
					sum -= x.(int64)
				}
				return sum
			}
		},
		"*": func(xs ...Any) Any {
			sum := int64(1)
			for _, x := range xs {
				sum *= x.(int64)
			}
			return sum
		},
		"<=": func(xs ...Any) Any {
			last := xs[0].(int64)
			for _, x := range xs {
				if x.(int64) < last {
					return false
				}
				last = x.(int64)
			}
			return true
		},
	})

	scanner := bufio.NewScanner(os.Stdin)
	for prompt(); scanner.Scan(); prompt() {
		in := scanner.Text()
		parsed, _ := read(in) // TODO handle err
		evaled := eval(parsed, env)
		fmt.Printf("%v (%T)\n", evaled, evaled)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
