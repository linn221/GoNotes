package utils

import "testing"

type Case struct {
	Input []int
	Test  int
	Want  bool
}

func TestExistsChecker(t *testing.T) {

	cases := []Case{
		{
			Input: []int{1, 2, 3},
			Test:  1,
			Want:  true,
		},
		{
			Input: []int{1, 2, 3},
			Test:  11,
			Want:  false,
		},
		{
			Input: []int{1, 2, 2, 3},
			Test:  2,
			Want:  true,
		},
		{
			Input: []int{1, 2, 2, 3},
			Test:  4,
			Want:  false,
		},
	}
	for _, cse := range cases {
		f := NewExistsChecker(cse.Input)
		got := f(cse.Test)
		if got != cse.Want {
			t.Errorf("input: %#v\ntest: %d\ngot: %v\n", cse.Input, cse.Test, got)
		}
	}
}
