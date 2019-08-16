package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func readRoutine(ctx context.Context, conn net.Conn) {
	scanner := bufio.NewScanner(conn)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() { // scan waites
				log.Printf("unable to scan")
				break OUTER
			}
			text := scanner.Text()
			log.Printf("From server: %s", text)
		}
	}
	log.Printf("finished write routine")
}

func writeRoutine(ctx context.Context, conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				break OUTER
			}
			str := scanner.Text()
			log.Printf("to server: %v\n", str)
			conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		}
	}
	log.Printf("finished write routine")
}

func main() {
	timeout := flag.Int("t", 3, "Connection timeout in seconds")
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatal("provide address and port")
	}
	address, port := flag.Args()[0], flag.Args()[1]

	log.Printf("go to %s:%s, wait %v seconds", address, port, *timeout)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(*timeout)*time.Second)

	d := net.Dialer{}
	conn, err := d.Dial("tcp", address+":"+port)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	fmt.Println(conn, err)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		readRoutine(ctx, conn)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		writeRoutine(ctx, conn)
		wg.Done()
	}()

	time.Sleep(12 * time.Minute)
	cancel()
	wg.Wait()
	conn.Close()
}
