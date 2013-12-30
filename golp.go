package main

import (
	"bufio"
	"fmt"
	"os"
)

type Object interface{}

type Symbol string

type Env struct {
	env map[Symbol]Object
}

func (e *Env) Find(s *Symbol) Object {
	if obj, ok := e.env[*s]; ok {
		return obj
	}
	return nil
}

func NewEnv() *Env {
	return &Env{}
}

func eval(sexp Object, env *Env) Object {
	switch sexp := sexp.(type) {
	case Symbol:
		return env.Find(&sexp)
	default:
		return sexp
	}
}

func prompt() {
	fmt.Print(">>> ")
}

func main() {
	//env := NewEnv()

	scanner := bufio.NewScanner(os.Stdin)
	for prompt(); scanner.Scan(); prompt() {
		in := scanner.Text()
		fmt.Println(in)
		//fmt.Println(eval(read, env))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
