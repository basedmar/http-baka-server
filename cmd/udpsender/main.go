package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	raddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("err from reading bytes from stream")

	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatalf("err from reading bytes from stream")

	}
	defer conn.Close()
	rdr := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">")
		output, err := rdr.ReadString('\n')
		if err != nil {
			log.Print(err.Error())
		}
		conn.Write([]byte(output))
	}
}
