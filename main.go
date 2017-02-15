package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/danielwhite/dumbledore/plugin"
)

var (
	addressFlag   = flag.String("tcp", "127.0.0.1:8080", "TCP service address")
	configFlag    = flag.String("f", "", "filter configuration file")
	pluginDirFlag = flag.String("plugin-dir", ".", "directory plugins can be loaded from")
)

var filters = []plugin.Filter{}

func main() {
	flag.Parse()

	if *configFlag == "" {
		fmt.Fprintf(os.Stderr, "configuration file must be specified\n")
		flag.Usage()
		os.Exit(1)
	}

	// Read configuration for the pipeline.
	cfg := readConfig(*configFlag)
	cfg.PluginDir = *pluginDirFlag

	// Load filters from plugins.
	filters = cfg.loadFilters()

	// Create a channel that processed events can be read from.
	out := startListener(*addressFlag)

	// Simply stream all output to the console.
	enc := json.NewEncoder(os.Stdout)
	for v := range out {
		enc.Encode(v)
	}
}

func startListener(addr string) <-chan plugin.Event {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to open listener: %s\n", err)
	}
	log.Printf("Listening on: %s\n", addr)

	out := make(chan plugin.Event)
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

func readAll(r io.Reader, out chan<- plugin.Event) error {
	in := make(chan plugin.Event)
	defer close(in)

	// Start a pipeline that filters events between in and
	// out. These are created per-connection.
	startFilters(filters, in, out)

	// For each event read from the connection, pass it to the
	// filter pipeline.
	dec := json.NewDecoder(r)
	for {
		var v plugin.Event
		if err := dec.Decode(&v); err != nil {
			return err
		}
		in <- v
	}
}

func startFilters(filters []plugin.Filter, in <-chan plugin.Event, out chan<- plugin.Event) {
	for _, filter := range filters {
		in = startFilter(filter, in)
	}
	go func() {
		for event := range in {
			out <- event
		}
	}()
}

func startFilter(filter plugin.Filter, in <-chan plugin.Event) <-chan plugin.Event {
	out := make(chan plugin.Event)
	go func() {
		for event := range in {
			filter.Filter(event, out)
		}
		close(out)
	}()
	return out
}
