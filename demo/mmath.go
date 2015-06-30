package mmath

import "log"

func Add(x, y int) int {
	return x + y
}

func Sub(x, y int) int {
	return x - y
}

func Mul(x, y int) int {
	return x * y
}

func Div(x, y int) int {
	return x / y
}

func Panic() {
    log.Panic()
}