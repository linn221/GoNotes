package formscanner

import (
	"cmp"
	"errors"
	"fmt"
	"time"
)

type ValidateFunc[T any] func(T) error

func MinMax(min int, max int) ValidateFunc[string] {
	return func(s string) error {
		if len(s) >= min && len(s) <= max {
			return nil
		}
		return fmt.Errorf("string length must be between %d and %d", min, max)
	}
}

func Gte[T cmp.Ordered](i T) ValidateFunc[T] {
	return func(v T) error {
		if v < i {
			return fmt.Errorf("must be greater than or equal %v", i)
		}
		return nil
	}
}

func Lte[T cmp.Ordered](i T) ValidateFunc[T] {
	return func(v T) error {
		if v > i {
			return fmt.Errorf("must be less than or equal %v", i)
		}
		return nil
	}
}

func Min(min int) ValidateFunc[string] {
	return func(s string) error {
		if len(s) >= min {
			return nil
		}
		return fmt.Errorf("string length must be greater than %d", min)
	}
}
func Max(max int) ValidateFunc[string] {
	return func(s string) error {
		if len(s) <= max {
			return nil
		}
		return fmt.Errorf("string length cannot be greater than %d", max)
	}
}

func InFuture(t time.Time) error {
	if t.IsZero() {
		return nil
	}
	if t.After(time.Now()) {
		return nil
	}

	return errors.New("time must be in the future")
}
