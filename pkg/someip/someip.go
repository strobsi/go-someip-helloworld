package someip

import (
	"encoding/binary"
	"errors"
)

func (header *SOMEIPHeader) Serialize() []byte {
	data := make([]byte, 16)

	binary.BigEndian.PutUint32(data[0:4], header.MessageID)
	binary.BigEndian.PutUint32(data[4:8], header.Length)
	data[8] = header.ProtocolVersion
	data[9] = header.InterfaceVersion
	data[10] = header.MessageType
	data[11] = header.ReturnCode

	return data
}

func ParseSOMEIPHeader(data []byte) (*SOMEIPHeader, error) {
	if len(data) < 16 {
		return nil, errors.New("insufficient data for SOME/IP header")
	}

	header := &SOMEIPHeader{
		MessageID:        binary.BigEndian.Uint32(data[0:4]),
		Length:           binary.BigEndian.Uint32(data[4:8]),
		RequestID:        binary.BigEndian.Uint32(data[8:12]),
		ProtocolVersion:  data[12],
		InterfaceVersion: data[13],
		MessageType:      data[14],
		ReturnCode:       data[15],
	}

	return header, nil
}

func NewSOMEIPMessage(messageID uint32, payload []byte) *SOMEIPMessage {
	return &SOMEIPMessage{
		Header: SOMEIPHeader{
			MessageID:        messageID,
			Length:           uint32(len(payload)),
			ProtocolVersion:  0x01,
			InterfaceVersion: 0x01,
			MessageType:      0x00, // 0x00: Request, 0x01: RequestNoReturn, 0x02: Notification, 0x80: Response, 0x81: Error
			ReturnCode:       0x00,
		},
		Payload: payload,
	}
}

func (msg *SOMEIPMessage) Serialize() []byte {
	data := make([]byte, 16+len(msg.Payload))

	binary.BigEndian.PutUint32(data[0:4], msg.Header.MessageID)
	binary.BigEndian.PutUint32(data[4:8], msg.Header.Length)
	binary.BigEndian.PutUint32(data[8:12], msg.Header.RequestID)
	data[12] = msg.Header.ProtocolVersion
	data[13] = msg.Header.InterfaceVersion
	data[14] = msg.Header.MessageType
	data[15] = msg.Header.ReturnCode

	copy(data[16:], msg.Payload)

	return data
}

func Deserialize(data []byte) (*SOMEIPMessage, error) {
	header, err := ParseSOMEIPHeader(data)
	if err != nil {
		return nil, err
	}

	payload := data[16:]

	return &SOMEIPMessage{
		Header:  *header,
		Payload: payload,
	}, nil
}

func (msg *SOMEIPSDMessage) Serialize() []byte {
	entryCount := len(msg.Entries)
	data := make([]byte, 16+32*entryCount)

	// Serialize the SOME/IP header
	headerData := msg.Header.Serialize()
	copy(data[:16], headerData)

	// Serialize the service entries
	for i, entry := range msg.Entries {
		offset := 16 + 32*i
		binary.BigEndian.PutUint16(data[offset:offset+2], entry.ServiceID)
		binary.BigEndian.PutUint16(data[offset+2:offset+4], entry.InstanceID)
		data[offset+4] = entry.MajorVersion
		binary.BigEndian.PutUint32(data[offset+8:offset+12], entry.MinorVersion)
		binary.BigEndian.PutUint32(data[offset+12:offset+16], entry.TTL)
		copy(data[offset+16:offset+20], entry.IPAddress.To4())
		binary.BigEndian.PutUint16(data[offset+20:offset+22], entry.Port)
		data[offset+21] = entry.Protocol
	}

	return data
}
