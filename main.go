package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device  string        = "wlo1"
	snapLen int           = 1024
	capNei  bool          = false
	timeout time.Duration = 30 * time.Second
)

func main() {
	handle, err := pcap.OpenLive(device, int32(snapLen), capNei, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	pclSrc := gopacket.NewPacketSource(handle, handle.LinkType())

	fmt.Println("start capturing...")
	for pck := range pclSrc.Packets() {
		fmt.Println("Handle packet!")
		ok := handlePacket(pck)
		if !ok {
			fmt.Println("some error")
		}
	}

}

// tcpdummi --p tcp --f write.txt
func handlePacket(pck gopacket.Packet) bool {
	if tcp := pck.Layer(layers.LayerTypeTCP); tcp != nil {
		tcpLay, ok := tcp.(*layers.TCP)
		if !ok {
			fmt.Println("Not TCP")
			return false
		}
		handleTCP(tcpLay)
	}
	return true

}

func handleTCP(tcp *layers.TCP) {
	fmt.Println("capture tcp")
	fmt.Println(tcp.SrcPort)
}
