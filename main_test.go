package main

import "testing"

func TestGradient(t *testing.T) {
	f := gradient(4, 1024)

	in, out := 0, "#FF0000"
	res := f(uint(in))
	t.Logf("f(%d) = %s", in, res)
	if res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 256, "#FF0000"
	res = f(uint(in))
	t.Logf("f(%d) = %s", in, res)
	if res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 257, "#FF7F00"
	res = f(uint(in))
	t.Logf("f(%d) = %s", in, res)
	if res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 512, "#FF7F00"
	res = f(uint(in))
	t.Logf("f(%d) = %s", in, res)
	if res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 513, "#7FFF00"
	res = f(uint(in))
	t.Logf("f(%d) = %s", in, res)
	if res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 768, "#7FFF00"
	res = f(uint(in))
	t.Logf("f(%d) = %s", in, res)
	if res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 769, "#00FF00"
	res = f(uint(in))
	t.Logf("f(%d) = %s", in, res)
	if res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 1024, "#00FF00"
	res = f(uint(in))
	t.Logf("f(%d) = %s", in, res)
	if res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}

	in, out = 1025, "#00FF00"
	res = f(uint(in))
	t.Logf("f(%d) = %s", in, res)
	if res != out {
		t.Errorf("f(%d) = %s, want %s", in, res, out)
	}
}
