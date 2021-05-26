package binframe

// HeadSize get varint head size, for performance, I use hard code instead of loops
func HeadSize(i uint64) int {

	if i < 1<<7 {
		return 1
	}
	if i < 1<<14 {
		return 2
	}
	if i < 1<<21 {
		return 3
	}
	if i < 1<<28 {
		return 4
	}
	if i < 1<<35 {
		return 5
	}
	if i < 1<<42 {
		return 6
	}
	if i < 1<<49 {
		return 7
	}
	if i < 1<<56 {
		return 8
	}
	if i < 1<<63 {
		return 9
	}

	return 10
}
