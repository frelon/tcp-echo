package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/golang/glog"
)

var port int

func init() {
	flag.IntVar(&port, "port", 7, "Port to listen to")
}

func handleConnection(connection net.Conn) {
	defer connection.Close()

	ch := make(chan []byte)
	errorCh := make(chan error)

	go func(ch chan []byte, errorCh chan error) {
		for {
			data := make([]byte, 1024)
			_, err := connection.Read(data)
			if err != nil {
				errorCh <- err
				return
			}

			ch <- data
		}
	}(ch, errorCh)

	for {
		select {
		case data := <-ch:
			connection.Write(data)
			break
		case err := <-errorCh:
			glog.Errorf("Error receiving: %v", err)
			break
		}
	}
}

func main() {
	flag.Parse()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	go func() {
		_ = <-sigc

		glog.Info("Received interrupt, closing")
		glog.Flush()
		os.Exit(1)
	}()

	glog.Infof("Starting up tcp-echo on port %+v", port)

	sock, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		glog.Fatalf("Error when opening socket: %+v", err)
	}

	glog.Info("Waiting for connections...")
	for {
		conn, err := sock.Accept()
		if err != nil {
			glog.Warningf("Error on accept: %+v", err)
		}

		glog.Infof("Connection accepted from %+v", conn.RemoteAddr())

		go handleConnection(conn)
	}
}
