package main

import (
	"fmt"
	"log"
	"net"
	"pro01/net_example/util"
	"time"
)

func main() {
	// Client Goroutine
	for i := 0; i < util.ClientNum; i++ {
		go func(i int) {
			connTo, err := net.Dial("udp", fmt.Sprintf("%s:%d", util.ServerIP, util.ServerPort))
			if err != nil {
				log.Fatal(err)
			}

			udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", util.ServerIP, util.ServerPort+i+1))
			if err != nil {
				log.Fatal(err)
			}
			connFrom, err := net.ListenUDP("udp", udpAddr)
			if err != nil {
				log.Fatal(err)
			}

			for {
				buf := []byte(fmt.Sprintf("hello from client%d...", i))
				connTo.Write(buf)

				buf0 := make([]byte, 100)
				connFrom.Read(buf0)
				log.Println(fmt.Sprintf("[Client%d]:%s", i, string(buf0)))
			}

		}(i)
	}

	// Server
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", util.ServerIP, util.ServerPort))
	if err != nil {
		log.Fatal(err)
	}
	connFrom, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	connTo := make([]net.Conn, util.ClientNum)
	for i := 0; i < util.ClientNum; i++ {
		conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", util.ServerIP, util.ServerPort+i+1))
		if err != nil {
			log.Fatal(err)
		}
		connTo[i] = conn
	}

	for {
		for i := 0; i < util.ClientNum; i++ {
			buf := make([]byte, 100)
			connFrom.Read(buf)
			log.Println("[Server]:", string(buf))

			buf0 := []byte("response from server...")
			connTo[i].Write(buf0)
		}
		time.Sleep(time.Second)
	}
}
