package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"
	"os"
)

var (
	addressFlag = flag.String("tcp", "127.0.0.1:8080", "TCP service address")
)

func main() {
	flag.Parse()

	out := startListener(*addressFlag)

	// Simply stream all output to the console.
	enc := json.NewEncoder(os.Stdout)
	for v := range out {
		enc.Encode(v)
	}
}

func startListener(addr string) <-chan map[string]interface{} {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to open listener: %s\n", err)
	}
	log.Printf("Listening on: %s\n", addr)

	out := make(chan map[string]interface{})
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Connect failed: %s\n", err)
			}

			go func() {
				log.Printf("Streaming input from %s", conn.RemoteAddr())

				// Read
				if err := readAll(conn, out); err == io.EOF {
					log.Printf("Connection from %s closed", conn.RemoteAddr())
				} else if err != nil {
					log.Printf("Error decoding input from %s: %s", conn.RemoteAddr(), err)
				}

				conn.Close()
			}()
		}
	}()
	return out
}

func readAll(r io.Reader, out chan<- map[string]interface{}) error {
	dec := json.NewDecoder(r)
	for {
		var v map[string]interface{}
		if err := dec.Decode(&v); err != nil {
			return err
		}
		out <- v
	}
}
