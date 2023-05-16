package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	ip     = "127.0.0.1"
	port   = 9090
	cliNum = 5
)

func getAddr(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

func getRecvUDPConn(ip string, port int) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", getAddr(ip, port))
	if err != nil {
		return nil, err
	}
	recvConn, err := net.ListenUDP("udp", udpAddr)
	return recvConn, err
}

func getSendUDPConn(ip string, port int) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", getAddr(ip, port))
	if err != nil {
		return nil, err
	}
	sendConn, err := net.DialUDP("udp", nil, udpAddr)
	return sendConn, err
}

func UDP() {
	// Client
	for i := 0; i < cliNum; i++ {
		go func(i int) {
			sendConn, err := getSendUDPConn(ip, port)
			if err != nil {
				log.Fatal(err)
			}
			defer sendConn.Close()

			recvConn, err := getRecvUDPConn(ip, port+i+1)
			if err != nil {
				log.Fatal(err)
			}
			defer recvConn.Close()

			for {
				sendConn.Write([]byte(fmt.Sprintf("hello from client%d...", i)))

				readBuf := make([]byte, 1024)
				recvConn.Read(readBuf)
				log.Printf("[Client%d]:%s", i, string(readBuf))
			}
		}(i)
	}

	// Server
	recvConn, err := getRecvUDPConn(ip, port)
	if err != nil {
		log.Fatal(err)
	}
	defer recvConn.Close()

	sendConns := make([]*net.UDPConn, cliNum)
	for i := 0; i < cliNum; i++ {
		sendConns[i], err = getSendUDPConn(ip, port+i+1)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer func() {
		for i := 0; i < cliNum; i++ {
			sendConns[i].Close()
		}
	}()

	for {
		for i := 0; i < cliNum; i++ {
			readBuf := make([]byte, 1024)
			recvConn.Read(readBuf)
			log.Printf("[Server]:%s", string(readBuf))

			sendConns[i].Write([]byte("hello from server..."))
		}
		time.Sleep(time.Second)
	}
}

func TCP() {
	// Client
	for i := 0; i < cliNum; i++ {
		go func(i int) {
			conn, err := net.Dial("tcp", getAddr(ip, port))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			for {
				conn.Write([]byte(fmt.Sprintf("hello from client%d...", i)))

				readBuf := make([]byte, 1024)
				conn.Read(readBuf)
				log.Printf("[Client%d]:%s", i, string(readBuf))
			}

		}(i)
	}

	// Server
	conns := make([]net.Conn, cliNum)

	listener, err := net.Listen("tcp", getAddr(ip, port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for i := 0; i < cliNum; i++ {
		conns[i], err = listener.Accept()
	}

	for {
		for i := 0; i < cliNum; i++ {
			readBuf := make([]byte, 1024)
			conns[i].Read(readBuf)
			log.Printf("[Server]:%s", string(readBuf))

			conns[i].Write([]byte("hello from server..."))
		}
		time.Sleep(time.Second)
	}
}

func main() {
	// UDP()
	TCP()
}
