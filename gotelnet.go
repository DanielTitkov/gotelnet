package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	timeout := flag.Int("t", 3, "Connection timeout in seconds")
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatal("Provide address and port")
	}
	address, port := flag.Args()[0], flag.Args()[1]

	log.Printf("go to %s:%s, wait %v seconds", address, port, *timeout)

	d := net.Dialer{}
	res, err := d.Dial("tcp", "yandex.com:80")
	fmt.Println(res, err)
}
