package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"math"
	"container/heap"

	"ping/v2/pkg/priorityQueue"

	"github.com/spf13/pflag"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	var count int
	pflag.IntVarP(&count, "count", "c", 5, "Number of messages to send (default: 5)")
	pflag.Parse()

	destination_IP := pflag.Args()

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./ping2.0v <DESTINATION_IP...> -c [MESSAGE_COUNT]")
		return
	}

	for i := 0; i < len(destination_IP); i++ {
		var receivedCnt int
		h := &priorityQueue.DurationHeap{}
		heap.Init(h)

		ips , err := net.LookupIP(destination_IP[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
			return
		}

		var seqNum int = 0	// Sequence number는 이후 자동으로 증가
		
		DataMessage := fmt.Sprintf("Ping for %s", destination_IP[i])
		fmt.Printf("PING %v (%v): %d data bytes\n", destination_IP[i], ips[0], len(DataMessage))

		for j := 0; j < count; j++ {
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
					Data: []byte(DataMessage), 	// 패킷 손실이나 데이터 변조 여부를 확인하기 위한 문자열
				},
			}

			packet, err := icmpMessage.Marshal(nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to marshal ICMP: %v\n", err)
				return
			}

			conn, err := net.Dial("ip4:icmp", destination_IP[i]) // arg[1]로 connect 시도
			if err != nil {
				fmt.Fprintf(os.Stderr, "Dial error: %v\n", err)
				return
			}
			defer conn.Close()

			if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
				fmt.Fprintf(os.Stderr, "Deadline error: %v\n", err)
				continue
			}

			start := time.Now()
			if _, err := conn.Write(packet); err != nil {
				fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
				return
			}

			reply := make([]byte, 1500)	// 대부분의 Ethernet에서는 MTU가 1500바이트이기 때문

			n, err := conn.Read(reply)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Request timeout for %s icmp_seq: %d\n", destination_IP[i], seqNum)
				seqNum++
				continue
			} else {
				receivedCnt++
			}

			duration := time.Since(start)
			parsedMessage, err := icmp.ParseMessage(1, reply[20:n]) // IP header가 20바이트 이므로
			if err != nil {
				fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
				return
			}


			switch parsedMessage.Type {
			case ipv4.ICMPTypeEchoReply:
				ipHeader, _ := ipv4.ParseHeader(reply[:20])
				ttl := ipHeader.TTL

				icmpSeq := parsedMessage.Body.(*icmp.Echo).Seq
				heap.Push(h, duration)
				fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%d time=%v\n", n, ips[0], icmpSeq, ttl, duration)
			default:
				fmt.Printf("Unexpected response: %+v\n", parsedMessage)
			}

			seqNum++
			time.Sleep(1 * time.Second)
		}

		fmt.Printf("\n--- %s ping statistics ---\n", destination_IP[i])
		lossRate := 100.0 - (float64(receivedCnt) / float64(seqNum)) * 100
		fmt.Printf("%d packets transmitted, %d packets received, %.1f%% packet loss\n", seqNum, receivedCnt, lossRate)
		
		// min / avg / max / stddev 구하기

		var durations []time.Duration

		for h.Len() > 0 {
			d := heap.Pop(h).(time.Duration)
			durations = append(durations, d)
		}

		min := durations[0]
		max := durations[len(durations) - 1]

		var total time.Duration
		for _, d := range durations {
			total += d
		}

		avg := time.Duration(int64(total) / int64(len(durations)))

		var varianceSum float64
		for _, d := range durations {
			diff := float64(d - avg)
			varianceSum += diff * diff
		}

		stddev := time.Duration(0)
		stddev = time.Duration(math.Sqrt(varianceSum / float64(len(durations))))

		fmt.Printf("round-trip min/avg/max/stddev = %.3f / %.3f / %.3f / %.3f ms\n\n", 
			float64(min.Microseconds()) / 1000,
			float64(avg.Microseconds()) / 1000,
			float64(max.Microseconds()) / 1000,
			float64(stddev.Microseconds()) / 1000,
		)
	}
}