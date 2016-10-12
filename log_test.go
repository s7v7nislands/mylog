package mylog

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"
)

func TestInit(t *testing.T) {
	out := &bytes.Buffer{}
	Init(INFO, out, 0, false)

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

func TestInitJSON(t *testing.T) {
	out := &bytes.Buffer{}
	Init(INFO, out, log.LstdFlags|log.Lshortfile, true)

	Debugf("debug")

	if out.String() != "" {
		t.Errorf("Debugf error: %s", out.String())
	}

	out.Reset()

	Infof("info")
	var m map[string]interface{}
	_ = json.Unmarshal(out.Bytes(), &m)
	message := m["msg"].(string)
	if message != "Info: info" {
		t.Errorf("Debugf error: %s", out.String())
	}

}
