package queue

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const bodyLen = 4

// binaryMessage describes binary message
type binaryMessage struct {
	Len  int
	Body []byte
}

func newBinaryMessage(data []byte) *binaryMessage {
	return &binaryMessage{
		Len:  len(data),
		Body: data,
	}
}

// serialize is to add log type & body length
func (b *binaryMessage) serialize() []byte {
	buf := make([]byte, len(b.Body)+bodyLen)
	binary.BigEndian.PutUint32(buf[:bodyLen], uint32(len(b.Body)))
	copy(buf[bodyLen:], b.Body)
	return buf
}

// deserialize turns bytes to binaryMessage
func deserialize(data []byte) (*binaryMessage, error) {
	// data size must be greater than 2 bytes
	if len(data) < bodyLen {
		return nil, fmt.Errorf("too short data size")
	}

	message := &binaryMessage{}
	message.Len = int(binary.BigEndian.Uint32(data[:bodyLen]))
	if (message.Len + bodyLen) != len(data) {
		return nil, fmt.Errorf("mismatch body length")
	}
	message.Body = data[bodyLen:]
	return message, nil
}

// deserializeSlice turns bytes to binaryMessage slice
func deserializeSlice(data []byte) ([]*binaryMessage, error) {
	c := bytes.NewReader(data)
	var ret []*binaryMessage

	for {
		// get the header
		header := make([]byte, bodyLen)
		if _, err := io.ReadFull(c, header); err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		message := &binaryMessage{}
		message.Len = int(binary.BigEndian.Uint32(header))

		// get the body
		bodyByte := make([]byte, message.Len)
		if _, err := io.ReadFull(c, bodyByte); err != nil {
			return nil, err
		}
		message.Body = bodyByte

		ret = append(ret, message)
	}

	return ret, nil
}
