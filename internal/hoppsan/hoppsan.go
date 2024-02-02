package hoppsan

import "testing"

func NoError(t testing.TB, err error) {
	if err == nil {
		return
	}
	t.Errorf("expected no error, got: %v\n", err)
	t.FailNow()
}

func Equal[T comparable](t testing.TB, left, right T) {
	if left == right {
		return
	}
	t.Errorf("expected left and right to be equal, got:\nLHS: %v\nRHS: %v\n", left, right)
}
