// Copyright 2018 visualfc. All rights reserved.

package tk

import "testing"

func init() {
	registerTest("Separator", testSeparator)
}

func testSeparator(t *testing.T) {
	w := NewSeparator(nil, Vertical)
	defer w.Destroy()

	w.SetOrient(Horizontal)
	if v := w.Orient(); v != Horizontal {
		t.Fatal("Orient", Horizontal, v)
	}

	w.SetTakeFocus(true)
	if v := w.IsTakeFocus(); v != true {
		t.Fatal("IsTakeFocus", true, v)
	}
}
