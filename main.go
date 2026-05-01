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
	/*
				  (linux)
				   создаёт структуру открытого файла, после возвращает сокет
							по этому сокету привязывает в заглушки структуры отркытого файла методы для работы с пакетами
					далее bind говорит к какому интерфейсу привязаться — записывает индекс eth0 в struct sock
					далее у ядра запрашивается новый ring buffer , в который оно будет писать сырые данные

				теперь handle держит fd + виртуальный адрес ring buffer
		                 ядро пишет пакеты туда через packet_rcv()
		                 код читает байты напрямую оттуда
	*/
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	pclSrc := gopacket.NewPacketSource(handle, handle.LinkType())
	// создается обертка над handle, которая знает какие байты отвечат за опред. уровни (tcp и т.п)

	fmt.Println("start capturing...")
	// Packets - создаёт канал, который сразу же возвращает, и создает горутинку, которая в бесконечном цикле читает сырые байты с помощью функции NextPacket
	// в свою очередь NextPacket спит до того момента, пока ядро не разбудит его из-за новых данных
	for pck := range pclSrc.Packets() {
		fmt.Println("Handle packet!")
		ok := parsePacket(pck)
		if !ok {
			fmt.Println("some error")
		}
	}

}

// tcpdummi --p tcp --f write.txt
func parsePacket(pck gopacket.Packet) bool {
	if tcp := pck.Layer(layers.LayerTypeTCP); tcp != nil {
		tcpLay, ok := tcp.(*layers.TCP) // type assertion
		if !ok {
			fmt.Println("Not TCP")
			return false
		}
	}
	return true

}

func handleTCP(tcp *layers.TCP) {
	fmt.Println("capture tcp")
	fmt.Println(tcp.SrcPort)
}
