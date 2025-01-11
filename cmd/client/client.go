package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	var port int
	var addr string

	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.StringVar(&addr, "address", "localhost", "--address localhost")
	flag.Parse()

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal("ResolveTCPAddr failed:", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal("Dial failed:", err)
	}

	go func() {
		scan := bufio.NewScanner(conn)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			fmt.Println(scan.Text())
		}
	}()

	for {
		inputScan := bufio.NewScanner(os.Stdin)
		inputScan.Split(bufio.ScanLines)
		for inputScan.Scan() {
			if inputScan.Text() == "exit" {
				return
			}
			conn.Write([]byte(fmt.Sprintf("%s\n", inputScan.Text())))
		}
	}
}
