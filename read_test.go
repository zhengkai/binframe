package binframe_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"sync"
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

	buf.Write([]byte{1, 19})
	msg, err = binframe.Read(buf)
	if err != nil || len(msg) != 1 || msg[0] != 19 {
		t.Error(`1 byte msg fail`)
	}
}

func TestChips(t *testing.T) {

	r, w := io.Pipe()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		msg, err := binframe.Read(r)
		if err != nil {
			t.Error(`read fail`, err)
		}

		bc := bytes.Compare(msg, []byte{2, 3, 1, 2, 3})
		if bc != 0 {
			t.Error(`msg mismatch`, bc)
		}
		wg.Done()
	}()

	w.Write([]byte{5, 2, 3})
	w.Write([]byte{1, 2, 3})
	wg.Wait()

	wg.Add(1)
	go func() {
		msg, err := binframe.Read(r)
		if err != nil {
			t.Error(`read fail`, err)
		}

		bc := bytes.Compare(msg, []byte{2, 3, 1, 2, 3})
		if bc != 0 {
			t.Error(`msg mismatch`, bc)
		}
		wg.Done()
	}()
	w.Write([]byte{5})
	w.Write([]byte{2, 3, 1, 2, 3})
	wg.Wait()
}

func TestMiddle(t *testing.T) {

	r, w := io.Pipe()
	var wg sync.WaitGroup

	randAB := make([]byte, 2000)
	rand.Read(randAB)

	buf := &bytes.Buffer{}
	binframe.Write(buf, randAB)

	content := buf.Bytes()

	wg.Add(1)
	go func() {
		msgA, err := binframe.Read(r)
		if err != nil {
			t.Error(`read fail`, err)
		}

		msgB, err := binframe.Read(r)
		if err != nil {
			t.Error(`read fail`, err)
		}

		bc := bytes.Compare(msgA, msgB)
		if bc != 0 {
			t.Error(`msg mismatch`, bc)
		}
		wg.Done()
	}()

	w.Write(content)
	w.Write(content[:300])
	w.Write(content[300:])
	wg.Wait()
}

func TestLong(t *testing.T) {

	r, w := io.Pipe()
	var wg sync.WaitGroup

	randAB := make([]byte, 3000000)
	rand.Read(randAB)

	buf := &bytes.Buffer{}
	binframe.Write(buf, randAB)

	content := buf.Bytes()

	wg.Add(1)
	go func() {

		// msg 1

		msgA, err := binframe.Read(r)
		if err != nil {
			t.Error(`read fail`, err)
		}

		// msg 2

		msgB, err := binframe.Read(r)
		if err != nil {
			t.Error(`read fail`, err)
		}

		bc := bytes.Compare(msgA, msgB)
		if bc != 0 {
			t.Error(`msg mismatch`, bc)
		}

		// msg 3 (close)

		_, err = binframe.Read(r)
		if err != io.ErrUnexpectedEOF {
			t.Error(`wrong err`, err)
		}

		wg.Done()
	}()

	w.Write(content)

	w.Write(content[:300])
	w.Write(content[300:])

	w.Write(content[:300])
	w.Write(content[300:10000])
	w.Close()

	wg.Wait()
}

func TestHead(t *testing.T) {

	r, w := io.Pipe()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		_, err := binframe.Read(r)
		if err != io.ErrUnexpectedEOF {
			t.Error(`read fail`, err)
		}
		wg.Done()
	}()

	w.Write([]byte{200, 200})
	w.Close()
	wg.Wait()
}
func TestShort(t *testing.T) {

	r, w := io.Pipe()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		_, err := binframe.Read(r)
		if err != io.ErrUnexpectedEOF {
			t.Error(`1 byte chip read fail`, err)
		}
		wg.Done()
	}()

	w.Write([]byte{1})
	w.Close()
	wg.Wait()

	buf := &bytes.Buffer{}
	buf.Write([]byte{1})
	_, err := binframe.Read(buf)
	if err != io.ErrUnexpectedEOF {
		t.Error(`read size with unexpected error`, err)
	}
}
