package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./ping1.0v <Destination>")
		return
	}
	// Target Address
	target := os.Args[1]

	var seqNum int = 1	// Sequence number는 이후 자동으로 증가
	for {
		// Message 만들기
		icmpMessage := icmp.Message {
			Type: ipv4.ICMPTypeEcho, // ICMPTypeEcho = 8
			Code: 0,

			// Checksum: 
			//	-> Marshal() 함수로 Checksum을 계산하여 바이트 배열에 포함시킴

			// Identifier, Seq Num, Payload(Data)는 icmp.Echo에 정의 되어 있음
			Body: &icmp.Echo {
				ID: os.Getpid() & 0xffff,	// ID는 주로 process id를 적는데, 16바이트로 맞추기 위해 0xffff로 자름
				Seq: seqNum,
				Data: []byte("FirstEcho"), 	// 패킷 손실이나 데이터 변조 여부를 확인하기 위한 문자열
			},
		}

		packet, err := icmpMessage.Marshal(nil)
		if err != nil {
			fmt.Println("Failed to marshal ICMP:", err)
			return
		}

		conn, err := net.Dial("ip4:icmp", target) // arg[1]로 connect 시도
		if err != nil {
			fmt.Println("Dial error:", err)
			return
		}
		defer conn.Close()

		if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
			fmt.Println("Deadline error:", err)
			continue
		}

		start := time.Now()
		if _, err := conn.Write(packet); err != nil {
			fmt.Println("Write error:", err)
			return
		}

		reply := make([]byte, 1500)	// 대부분의 Ethernet에서는 MTU가 1500바이트이기 때문

		n, err := conn.Read(reply)
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		duration := time.Since(start)
		parsedMessage, err := icmp.ParseMessage(1, reply[20:n]) // IP header가 20바이트 이므로
		if err != nil {
			fmt.Println("Parse error:", err)
			return
		}

		switch parsedMessage.Type {
		case ipv4.ICMPTypeEchoReply:
			fmt.Printf("Reply from %s: time=%v\n", target, duration)
		default:
			fmt.Printf("Unexpected response: %+v\n", parsedMessage)
		}

		seqNum++
		time.Sleep(1 * time.Second)
	}
}