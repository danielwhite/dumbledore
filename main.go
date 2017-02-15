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

var filters = []Filter{
	&pruneFilter{BlacklistNames: []string{"secret"}},
	&cloneFilter{Clones: []string{"Mini Me"}},
}

func main() {
	flag.Parse()

	out := startListener(*addressFlag)

	// Simply stream all output to the console.
	enc := json.NewEncoder(os.Stdout)
	for v := range out {
		enc.Encode(v)
	}
}

func startListener(addr string) <-chan Event {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to open listener: %s\n", err)
	}
	log.Printf("Listening on: %s\n", addr)

	out := make(chan Event)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Connect failed: %s\n", err)
			}

			// For each connection, read JSON from in
			// stream, and write the filtered results to
			// the output channel.
			go func() {
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

func readAll(r io.Reader, out chan<- Event) error {
	in := make(chan Event)
	defer close(in)

	// Start a pipeline that filters events between in and
	// out. These are created per-connection.
	startFilters(filters, in, out)

	// For each event read from the connection, pass it to the
	// filter pipeline.
	dec := json.NewDecoder(r)
	for {
		var v Event
		if err := dec.Decode(&v); err != nil {
			return err
		}
		in <- v
	}
}

func startFilters(filters []Filter, in <-chan Event, out chan<- Event) {
	for _, filter := range filters {
		in = startFilter(filter, in)
	}
	go func() {
		for event := range in {
			out <- event
		}
		log.Print("closing channel for filter set")
	}()
}
