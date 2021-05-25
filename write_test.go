package binframe_test

import (
	"io"
	"testing"

	"github.com/zhengkai/binframe"
)

func TestWriteFail(t *testing.T) {

	r, w := io.Pipe()
	r.Close()

	err := binframe.Write(w, []byte{1, 2, 3})
	if err == nil {
		t.Error(`write fail but no error`)
	}
}
