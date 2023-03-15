// server.go

package main

import (
	"fmt"
	"net"

	"github.com/strobsi/someip-go/pkg/config"
	"github.com/strobsi/someip-go/pkg/hello_world"
	"github.com/strobsi/someip-go/pkg/someip"
)

func main() {
	provider := config.NewE01HelloWorldProvider()
	method := config.NewE01HelloWorldMethod()

	addr := fmt.Sprintf("%s:%d", provider.SomeIpUnicastAddress, provider.SomeIpReliableUnicastPort)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Listening on %s...\n", addr)

	buffer := make([]byte, 4096)

	for {
		n, src, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			continue
		}

		request, err := someip.Deserialize(buffer[:n])
		if err != nil {
			fmt.Println("Error deserializing SOME/IP message:", err)
			continue
		}

		messageID := (uint32(provider.SomeIpInstanceID) << 16) | uint32(method.SomeIpMethodID)

		if request.Header.MessageID == messageID {
			fmt.Println("Received SayHello request")

			// Deserialize the request payload
			requestPayload, err := hello_world.DeserializeSayHelloRequest(request.Payload)
			if err != nil {
				fmt.Printf("Error deserializing request payload: %v\n", err)
				continue
			}

			responsePayload := fmt.Sprintf("Hello, %s!", requestPayload)

			// Serialize the response payload
			responseData, err := hello_world.SerializeSayHelloResponse(responsePayload)
			if err != nil {
				fmt.Printf("Error serializing response payload: %v\n", err)
				continue
			}

			response := someip.NewSOMEIPMessage(messageID, responseData)
			response.Header.MessageType = 0x80

			_, err = conn.WriteTo(response.Serialize(), src)
			if err != nil {
				fmt.Println("Error sending SayHello response:", err)
				continue
			}

			fmt.Println("SayHello response sent")
		} else {
			fmt.Println("Received unexpected message")
		}
	}
}
