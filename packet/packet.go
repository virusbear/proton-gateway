package packet

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"time"
)

type Packet interface {
	Mac() string
	Timestamp() time.Time
	Payload() []byte
}

type packetImpl struct {
	mac       string
	timestamp time.Time
	payload   []byte
}

func (packet packetImpl) Mac() string {
	return packet.mac
}

func (packet packetImpl) Timestamp() time.Time {
	return packet.timestamp
}

func (packet packetImpl) Payload() []byte {
	return packet.payload
}

func Read(reader io.Reader) (Packet, error) {
	mac, err := readHexString(reader, 6)
	if err != nil {
		return nil, err
	}
	var dataLen uint8
	err = binary.Read(reader, binary.LittleEndian, &dataLen)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, dataLen)
	_, err = reader.Read(payload)
	return packetImpl{
		mac:       *mac,
		timestamp: time.Now(),
		payload:   payload,
	}, err
}

func readHexString(reader io.Reader, length int) (*string, error) {
	bytes := make([]byte, length)
	read := 0
	for read < length {
		n, err := reader.Read(bytes[read:])
		if err != nil {
			return nil, err
		}
		read += n
	}
	// _, err := reader.Read(bytes)
	// if err != nil {
	// 		return nil, err
	// }
	value := hex.EncodeToString(bytes)
	return &value, nil
}
