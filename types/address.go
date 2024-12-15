package types

import (
	"encoding/hex"
	"fmt"
)

type Address [20]uint8

func (a Address) ToSlice() []byte {
	b := make([]byte, len(a))

	for i, v := range a {
		b[i] = byte(v)
	}

	return b
}

func AddressFromBytes(b []byte) Address {

	if len(b) != 20 {
		msg := fmt.Sprintf("given bytes with length %d should be 20 bytes ", len(b))
		panic(msg)
	}

	var value [20]uint8

	for i := 0; i < 20; i++ {
		value[i] = uint8(b[i])
	}

	return Address(value)
}

func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice())
}
