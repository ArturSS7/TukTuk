package dnslistener

import (
	"database/sql"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

type DnsMsg struct {
	Timestamp       string
	SourceIP        string
	DestinationIP   string
	DnsQuery        string
	DnsAnswer       []string
	DnsAnswerTTL    []string
	NumberOfAnswers string
	DnsResponseCode string
	DnsOpCode       string
}

func StartDNS(db *sql.DB) {

	device := "en0"

	fmt.Println("Capturing dns packets from interface " + device)

	handle, err := pcap.OpenLive(device, 1024, false, -1*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		dnsLayer := packet.Layer(layers.LayerTypeDNS)
		if dnsLayer != nil {
			ipLayer4 := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer4 != nil {
				ip, _ := ipLayer4.(*layers.IPv4)
				fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
			}
			ipLayer6 := packet.Layer(layers.LayerTypeIPv6)
			if ipLayer6 != nil {
				ip, _ := ipLayer6.(*layers.IPv6)
				fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
			}
		}
	}
}
