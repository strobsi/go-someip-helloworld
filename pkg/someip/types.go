package someip

import "net"

// SOMEIPHeader
type SOMEIPHeader struct {
	MessageID        uint32
	Length           uint32
	RequestID        uint32
	ProtocolVersion  uint8
	InterfaceVersion uint8
	MessageType      uint8
	ReturnCode       uint8
}

// SOMEIP Message
type SOMEIPMessage struct {
	Header  SOMEIPHeader
	Payload []byte
}

// Add this structure to represent SOME/IP-SD messages
type SOMEIPSDMessage struct {
	Header  SOMEIPHeader
	Entries []SOMEIPServiceEntry
}

type SOMEIPServiceEntry struct {
	ServiceID    uint16
	InstanceID   uint16
	MajorVersion uint8
	MinorVersion uint32
	TTL          uint32
	IPAddress    net.IP
	Port         uint16
	Protocol     uint8 // 0x06 for TCP, 0x11 for UDP
}
