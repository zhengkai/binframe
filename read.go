package binframe

import (
	"encoding/binary"
	"io"
)

// Read ...
func Read(r io.Reader) (msg []byte, err error) {

	size, body, err := getSize(r)
	if err != nil {
		return
	}

	if size == 1 {
		msg = []byte{body}
		return
	}

	is := int(size)

	msg = make([]byte, is)
	if is < 128 {
		msg[0] = body
		err = getBody(r, is-1, msg[1:])
	} else {
		err = getBody(r, is, msg)
	}

	return
}

func getBody(reader io.Reader, size int, msg []byte) (err error) {

	n, err := reader.Read(msg)
	if n == size {
		return
	}

	offset := n
	for {

		n, err = reader.Read(msg[offset:])
		if n == 0 {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			break
		}

		offset += n
		if offset == size {
			break
		}
	}

	return
}

func getSize(reader io.Reader) (size uint64, body byte, err error) {

	head := make([]byte, 2)

	var n int

	n, err = reader.Read(head)
	if n == 1 {
		if err != nil {
			return
		}
		n, err = reader.Read(head[1:])
		if n != 1 {
			if err == nil {
				err = io.ErrUnexpectedEOF
			}
			return
		}
	}

	if head[0] < 128 {
		size = uint64(head[0])
		body = head[1]
		return
	}

	for {
		size, n = binary.Uvarint(head)
		if n > 0 {
			break
		}
		if n < 0 {
			err = ErrSizeOverflow
			return
		}

		c := []byte{0}
		n, err = reader.Read(c)
		if err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return
		}
		head = append(head, c[0])
	}

	return
}
