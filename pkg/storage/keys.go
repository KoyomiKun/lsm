package storage

type ByteSliceKey []byte

func (k ByteSliceKey) Less(k2 ByteSliceKey) int {
	for i, b := range k {
		if i >= len(k2) {
			return -1
		}
		if b < k2[i] {
			return 1
		}
		if b > k2[i] {
			return -1
		}
	}
	if len(k) == len(k2) {
		return 0
	}
	return 1
}

type IntKey int

func (k IntKey) Less(k2 interface{}) int {
	k2 = k2.()
	if k < k2 {
		return 1
	}

	if k == k2 {
		return 0
	}

	return -1
}
