// client.go

package main

import (
	"fmt"
	"net"
	"time"

	"github.com/strobsi/someip-go/pkg/config"
	"github.com/strobsi/someip-go/pkg/hello_world"
	"github.com/strobsi/someip-go/pkg/someip"
)

func main() {
	provider := config.NewE01HelloWorldProvider()
	method := config.NewE01HelloWorldMethod()

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", provider.SomeIpUnicastAddress, provider.SomeIpReliableUnicastPort))
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	messageID := (uint32(provider.SomeIpInstanceID) << 16) | uint32(method.SomeIpMethodID)
	payload := "John Doe"

	// Serialize the payload
	requestPayload, err := hello_world.SerializeSayHelloRequest(payload)
	if err != nil {
		fmt.Printf("Error serializing payload: %v\n", err)
		return
	}

	msg := someip.NewSOMEIPMessage(messageID, requestPayload)
	_, err = conn.Write(msg.Serialize())
	if err != nil {
		fmt.Println("Error sending SayHello message:", err)
		return
	}

	fmt.Println("SayHello message sent")

	buffer := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving SOME/IP response:", err)
		return
	}

	response, err := someip.Deserialize(buffer[:n])
	if err != nil {
		fmt.Println("Error deserializing SOME/IP response:", err)
		return
	}

	if response.Header.MessageID == messageID && response.Header.MessageType == 0x80 {
		fmt.Println("Received SayHello response")

		// Deserialize the response payload
		responsePayload, err := hello_world.DeserializeSayHelloResponse(response.Payload)
		if err != nil {
			fmt.Printf("Error deserializing response payload: %v\n", err)
			return
		}

		fmt.Printf("Received message: %s\n", responsePayload)
	} else {
		fmt.Println("Received unexpected message")
	}
}
