package binframe_test

import (
	"encoding/binary"
	"testing"

	"github.com/zhengkai/binframe"
)

func TestHeadSize(t *testing.T) {

	buf := make([]byte, binary.MaxVarintLen64)

	for i := 7; i <= 64; i += 7 {
		for j := uint64(0); j <= 1; j++ {

			num := uint64(1<<i) - j

			n := binary.PutUvarint(buf, num)

			if n != binframe.HeadSize(num) {
				t.Error(`size mismatch with length`, num)
			}
		}
	}
}
