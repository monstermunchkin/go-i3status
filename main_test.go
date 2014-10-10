package main

import "testing"

func TestGradient(t *testing.T) {
	f := gradient(0, 3, 1024)

	in, out := 0, "#FF0000"
	if res := f(uint(in)); res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 256, "#FF0000"
	if res := f(uint(in)); res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 257, "#FFAA00"
	if res := f(uint(in)); res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 512, "#FFAA00"
	if res := f(uint(in)); res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 513, "#AAFF00"
	if res := f(uint(in)); res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 768, "#AAFF00"
	if res := f(uint(in)); res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 769, "#00FF00"
	if res := f(uint(in)); res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 1024, "#00FF00"
	if res := f(uint(in)); res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 1025, "#00FF00"
	if res := f(uint(in)); res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}
}
