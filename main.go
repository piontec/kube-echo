package main

import (
	"log"
	"net"
	"os"
)

func handleUDPConnection(conn *net.UDPConn, logEnabled bool) {
	buffer := make([]byte, 1024)

	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Println("Error reading UDP message:", err)
		return
	}

	message := string(buffer[:n])
	if logEnabled {
		log.Printf("UDP received from %s: %s\n", addr.String(), message)
	}

	_, err = conn.WriteToUDP(buffer[:n], addr)
	if err != nil {
		log.Println("Error sending UDP response:", err)
		return
	}
}

func handleTCPConnection(conn net.Conn, logEnabled bool) {
	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		log.Println("Error reading TCP message:", err)
		return
	}

	message := string(buffer[:n])
	if logEnabled {
		log.Printf("TCP received from %s: %s\n", conn.RemoteAddr().String(), message)
	}

	_, err = conn.Write(buffer[:n])
	if err != nil {
		log.Println("Error sending TCP response:", err)
		conn.Close()
		return
	}

	conn.Close()
}

func main() {
	listeningPort := os.Getenv("LISTEN_PORT")
	if listeningPort == "" {
		listeningPort = "7777"
	}

	logEnabled := os.Getenv("LOG_ENABLED") == "true"

	// UDP
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:"+listeningPort)
	if err != nil {
		log.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println("Error listening on UDP:", err)
		os.Exit(1)
	}

	defer udpConn.Close()

	log.Println("UDP Echo Server listening on", udpAddr.String())

	go func() {
		for {
			handleUDPConnection(udpConn, logEnabled)
		}
	}()

	// TCP
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:"+listeningPort)
	if err != nil {
		log.Println("Error resolving TCP address:", err)
		os.Exit(1)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println("Error listening on TCP:", err)
		os.Exit(1)
	}

	defer tcpListener.Close()

	log.Println("TCP Echo Server listening on", tcpAddr.String())

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Println("Error accepting TCP connection:", err)
			continue
		}

		go handleTCPConnection(conn, logEnabled)
	}
}
