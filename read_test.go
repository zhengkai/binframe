package binframe_test

import (
	"bytes"
	"testing"

	"github.com/zhengkai/binframe"
)

func TestBasic(t *testing.T) {

	content := []byte{1, 2, 3}

	buf := &bytes.Buffer{}

	err := binframe.Write(buf, content)
	if err != nil {
		t.Error(`write fail`, err)
	}

	msg, err := binframe.Read(buf)
	if err != nil {
		t.Error(`read fail`, err)
	}

	bc := bytes.Compare(msg, content)
	if bc != 0 {
		t.Error(`msg mismatch`, bc)
	}
}
