// Package fizzbuzz implements the generalized FizzBuzz algorithm used by the API.
package fizzbuzz

import (
	"errors"
	"strconv"
)

type Request struct {
	Int1  int
	Int2  int
	Limit int
	Str1  string
	Str2  string
}

const MaxLimit = 1_000_000

var (
	ErrInt1NotPositive  = errors.New("int1 must be a positive integer")
	ErrInt2NotPositive  = errors.New("int2 must be a positive integer")
	ErrLimitNotPositive = errors.New("limit must be a positive integer")
	ErrLimitTooLarge    = errors.New("limit exceeds the maximum allowed value")
	ErrStr1Empty        = errors.New("str1 must not be empty")
	ErrStr2Empty        = errors.New("str2 must not be empty")
)

func (r Request) Validate() error {
	switch {
	case r.Int1 <= 0:
		return ErrInt1NotPositive
	case r.Int2 <= 0:
		return ErrInt2NotPositive
	case r.Limit <= 0:
		return ErrLimitNotPositive
	case r.Limit > MaxLimit:
		return ErrLimitTooLarge
	case r.Str1 == "":
		return ErrStr1Empty
	case r.Str2 == "":
		return ErrStr2Empty
	}
	return nil
}

func GenerateFizzbuzz(r Request) []string {
	out := make([]string, r.Limit)
	for i := 1; i <= r.Limit; i++ {
		switch {
		case i%r.Int1 == 0 && i%r.Int2 == 0:
			out[i-1] = r.Str1 + r.Str2
		case i%r.Int1 == 0:
			out[i-1] = r.Str1
		case i%r.Int2 == 0:
			out[i-1] = r.Str2
		default:
			out[i-1] = strconv.Itoa(i)
		}
	}
	return out
}
