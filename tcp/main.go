package main

import (
	"fmt"
	"log"
	"net"
	"pro01/net_example/util"
	"time"
)

func main() {

	// Server Goroutine
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", util.ServerIP, util.ServerPort))
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()
		conns := make([]net.Conn, util.ClientNum)
		for i := 0; i < util.ClientNum; i++ {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			conns[i] = conn
		}
		for {
			for i := 0; i < util.ClientNum; i++ {
				buf := make([]byte, 100)
				conns[i].Read(buf)
				//fmt.Println(i)
				log.Println("[Server]:", string(buf))

				buf0 := []byte("response from server...")
				conns[i].Write(buf0)
			}
		}
	}()

	// Client
	conns := make([]net.Conn, util.ClientNum)
	for i := 0; i < util.ClientNum; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", util.ServerIP, util.ServerPort))
		if err != nil {
			log.Fatal(err)
		}
		conns[i] = conn
	}
	for {
		for i := 0; i < util.ClientNum; i++ {
			buf0 := []byte(fmt.Sprintf("hello from client%d...", i))
			conns[i].Write(buf0)

			buf := make([]byte, 100)
			conns[i].Read(buf)
			log.Println(fmt.Sprintf("[Client%d]:", i), string(buf))
		}
		time.Sleep(time.Second)
	}

}
