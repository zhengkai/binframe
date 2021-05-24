package binframe

import (
	"encoding/binary"
	"io"
)

// Write message
func Write(w io.Writer, msg []byte) (err error) {

	buf := make([]byte, binary.MaxVarintLen64)

	n := binary.PutUvarint(buf, uint64(len(msg)))
	_, err = w.Write(buf[:n])
	if err != nil {
		return
	}

	_, err = w.Write(msg)

	return
}
