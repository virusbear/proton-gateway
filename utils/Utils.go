package utils

import (
	"encoding/hex"
	"io"
)

func ReadMac(reader io.Reader) (*string, error) {
	bytes := make([]byte, 6)
	read := 0
	for read < 6 {
		n, err := reader.Read(bytes[read:])
		if err != nil {
			return nil, err
		}
		read += n
	}
	mac := hex.EncodeToString(bytes)
	return &mac, nil
}
