package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

var (
	Address = flag.String("tcp", "127.0.0.1:8080", "TCP service address")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", *Address)
	if err != nil {
		log.Fatalf("Failed to open listener: %s\n", err)
	}

	log.Printf("Listening on: %s\n", *Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connect failed: %s\n", err)
		}

		file, err := ioutil.TempFile(os.TempDir(), "dumbledore")
		if err != nil {
			log.Printf("Failed to open file: %s\n", err)
		}

		log.Printf("Streaming input from %s to %s", conn.RemoteAddr(), file.Name())

		go func() {
			defer file.Close()
			defer conn.Close()

			pipe := &Pipe{source: conn, dest: file}
			pipe.transfer()
		}()
	}
}

type Pipe struct {
	source io.Reader
	dest   io.WriteCloser
}

func (p *Pipe) transfer() {
	defer p.dest.Close()

	if _, err := io.Copy(p.dest, p.source); err != nil {
		log.Printf("Copy failed: %s\n", err)
	}
}
