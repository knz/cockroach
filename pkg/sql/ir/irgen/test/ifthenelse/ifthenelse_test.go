package test

import (
	"fmt"
	"testing"
)

func TestOptField(t *testing.T) {
	var c AllocContext
	stmt := c.If(123, nil, nil)

	if stmt.Then() != nil {
		t.Fatal("empty is not nil")
	}
}

func TestIfThenElseFmt(t *testing.T) {
	var c AllocContext

	testdata := []struct {
		stmt   *If
		expSQL string
	}{
		{c.If(123, c.Basic(), nil),
			`IF 123 THEN BASIC END`},
		{c.If(123, c.Basic(), c.Basic()),
			`IF 123 THEN BASIC ELSE BASIC END`},
		{c.If(123, c.If(456, c.Basic(), c.Basic()), c.Basic()),
			`IF 123 THEN IF 456 THEN BASIC ELSE BASIC END ELSE BASIC END`},
	}

	for i, d := range testdata {
		fmt.Println(d.stmt.String())

		sql := d.stmt.SQL()
		if sql != d.expSQL {
			t.Fatalf("%d: expected: %q, got: %q", i, sql, d.expSQL)
		}
	}
}
