package ps

import (
	"fmt"
	"testing"
)

func TestPSOutput1(t *testing.T) {
	input := "doe            51061   0.0  0.2  4391824  36344   ??  S    14Dec18   2:58.84 /Users/foo/go/bin/gocode -s -sock unix -addr 127.0.0.1:37373"
	proc, err := parsePSOutput(input)
	if err != nil {
		t.Errorf("err got %v want nil", err)
	}
	compare := func(a, b interface{}) {
		if a != b {
			t.Errorf("got %v want %v", a, b)
		}
	}
	compare(proc.USER, "doe")
	compare(proc.PID, int64(51061))
	compare(proc.CPU, 0.0)
	compare(proc.MEM, 0.2)
	compare(proc.VSZ, int64(4391824))
	compare(proc.RSS, int64(36344))
	compare(proc.TTY, "??")
	compare(proc.STAT, "S")
	compare(proc.START, "14Dec18")
	compare(proc.TIME, "2:58.84")
	compare(proc.COMMAND, "/Users/foo/go/bin/gocode -s -sock unix -addr 127.0.0.1:37373")
}

func TestPSOutput2(t *testing.T) {
	input := "root      1234  0.0  0.7 887956  7880 ?        Sl    2018   0:09 ./foo --arg=2"
	proc, err := parsePSOutput(input)
	if err != nil {
		t.Errorf("err got %v want nil", err)
	}
	compare := func(a, b interface{}) {
		if a != b {
			fmt.Printf("%T %T\n", a, b)
			t.Errorf("got %v want %v", a, b)
		}
	}
	compare(proc.USER, "root")
	compare(proc.PID, int64(1234))
	compare(proc.CPU, 0.0)
	compare(proc.MEM, 0.7)
	compare(proc.VSZ, int64(887956))
	compare(proc.RSS, int64(7880))
	compare(proc.TTY, "?")
	compare(proc.STAT, "Sl")
	compare(proc.START, "2018")
	compare(proc.TIME, "0:09")
	compare(proc.COMMAND, "./foo --arg=2")
}

func TestMultipleProcs(t *testing.T) {
	input := `doe            51061   0.0  0.2  4391824  36344   ??  S    14Dec18   2:58.84 /Users/foo/go/bin/gocode -s -sock unix -addr 127.0.0.1:37373
                  doe            51061   0.0  0.2  4391824  36344   ??  S    14Dec18   2:58.84 /Users/foo/go/bin/gocode -s -sock unix -addr 127.0.0.1:37373`
	proc, err := parsePSOutput(input)
	if err != ErrMultipleProcs {
		t.Errorf("err got %v want %v", err, ErrMultipleProcs)
	}
	if proc != nil {
		t.Errorf("proc got %v want nil", proc)
	}
}

func TestBadPSOutput(t *testing.T) {
	input := "a b c d e"
	proc, err := parsePSOutput(input)
	if err != ErrParsingPSOutput {
		t.Errorf("err got %v want %s", err, ErrParsingPSOutput)
	}
	if proc != nil {
		t.Errorf("proc got %v want nil", proc)
	}
}

func TestEmptyPSOutput(t *testing.T) {
	input := ""
	proc, err := parsePSOutput(input)
	if err != ErrParsingPSOutput {
		t.Errorf("err got %v want %s", err, ErrParsingPSOutput)
	}
	if proc != nil {
		t.Errorf("proc got %v want nil", proc)
	}
}
