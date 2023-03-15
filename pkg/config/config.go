package config

import (
	"net"
)

const (
	ServiceID             uint16 = 4660
	MethodID              uint32 = 30000
	Reliable              bool   = true
	StringEncoding        string = "utf16le"
	InstanceID            string = "commonapi.examples.HelloWorld"
	SomeIPInstanceID      uint16 = 22136
	UnicastAddress               = "192.168.178.108"
	ReliableUnicastPort          = 30499
	UnreliableUnicastPort        = 30499
)

type E01HelloWorldMethod struct {
	SomeIpMethodID       uint32
	SomeIpReliable       bool
	SomeIpStringEncoding string
}

type E01HelloWorldProvider struct {
	InstanceID                  string
	SomeIpInstanceID            uint16
	SomeIpUnicastAddress        net.IP
	SomeIpReliableUnicastPort   int
	SomeIpUnreliableUnicastPort int
}

func NewE01HelloWorldProvider() *E01HelloWorldProvider {
	return &E01HelloWorldProvider{
		InstanceID:                  InstanceID,
		SomeIpInstanceID:            SomeIPInstanceID,
		SomeIpUnicastAddress:        net.ParseIP(UnicastAddress),
		SomeIpReliableUnicastPort:   ReliableUnicastPort,
		SomeIpUnreliableUnicastPort: UnreliableUnicastPort,
	}
}

func NewE01HelloWorldMethod() *E01HelloWorldMethod {
	return &E01HelloWorldMethod{
		SomeIpMethodID:       MethodID,
		SomeIpReliable:       Reliable,
		SomeIpStringEncoding: StringEncoding,
	}
}
