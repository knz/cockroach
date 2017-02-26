package test

import (
	"fmt"
	"testing"
)

func TestSimple(t *testing.T) {
	var ctx AllocContext
	a := ctx.A()
	b := ctx.B()
	ca := ctx.C(a)
	cb := ctx.C(b)

	var field T

	// Retrieve the field "c" and check it implements the interface T.
	field = ca.C()
	// Also check the result is actually the original A in disguise.
	if p, ok := field.(*A); !ok || p != a {
		t.Fatalf("can't retrieve A")
	}

	// Same for B.
	field = cb.C()
	if p, ok := field.(*B); !ok || p != b {
		t.Fatalf("can't retrieve B")
	}

	// Try changing it!
	ca.SetC(&ctx, field)
	if p, ok := ca.C().(*B); !ok || p != b {
		t.Fatalf("can't retrieve B")
	}

}

func TestSimpleMarshal(t *testing.T) {
	var ctx AllocContext

	orig := ctx.C(ctx.B())
	ser, err := orig.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("serialized: %q\n", ser)

	var next C
	if err := next.Unmarshal(ser); err != nil {
		t.Fatal(err)
	}

	if _, ok := next.C().(*B); !ok {
		t.Fatal("can't convert back to B")
	}
}
