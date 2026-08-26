package fizzbuzz

import (
	"errors"
	"strings"
	"testing"
)

func TestGenerateFizzbuzz_Classic(t *testing.T) {
	got := GenerateFizzbuzz(Request{Int1: 3, Int2: 5, Limit: 16, Str1: "fizz", Str2: "buzz"})
	want := strings.Split("1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16", ",")

	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestGenerateFizzbuzz_CustomParameters(t *testing.T) {
	got := GenerateFizzbuzz(Request{Int1: 2, Int2: 7, Limit: 14, Str1: "foo", Str2: "bar"})
	want := []string{"1", "foo", "3", "foo", "5", "foo", "bar", "foo", "9", "foo", "11", "foo", "13", "foobar"}

	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestGenerateFizzbuzz_LimitZeroOrNegativeYieldsEmpty(t *testing.T) {
	if got := GenerateFizzbuzz(Request{Int1: 3, Int2: 5, Limit: 0, Str1: "a", Str2: "b"}); len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		req  Request
		want error
	}{
		{"valid", Request{3, 5, 100, "fizz", "buzz"}, nil},
		{"int1 zero", Request{0, 5, 100, "fizz", "buzz"}, ErrInt1NotPositive},
		{"int1 negative", Request{-1, 5, 100, "fizz", "buzz"}, ErrInt1NotPositive},
		{"int2 zero", Request{3, 0, 100, "fizz", "buzz"}, ErrInt2NotPositive},
		{"limit zero", Request{3, 5, 0, "fizz", "buzz"}, ErrLimitNotPositive},
		{"limit negative", Request{3, 5, -5, "fizz", "buzz"}, ErrLimitNotPositive},
		{"limit too large", Request{3, 5, MaxLimit + 1, "fizz", "buzz"}, ErrLimitTooLarge},
		{"str1 empty", Request{3, 5, 100, "", "buzz"}, ErrStr1Empty},
		{"str2 empty", Request{3, 5, 100, "fizz", ""}, ErrStr2Empty},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.req.Validate(); !errors.Is(err, tt.want) {
				t.Errorf("got %v, want %v", err, tt.want)
			}
		})
	}
}
