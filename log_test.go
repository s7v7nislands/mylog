package mylog

import (
	"bytes"
	"testing"
)

func TestInit(t *testing.T) {
	out := &bytes.Buffer{}
	Init(INFO, out, 0)

	Debugf("debug")

	if out.String() != "" {
		t.Errorf("Debugf error: %s", out.String())
	}

	out.Reset()

	Infof("info")
	if out.String() != "Info: info\n" {
		t.Errorf("Debugf error: %s", out.String())
	}

}
