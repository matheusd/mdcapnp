// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

var flagPort = flag.Int("port", 0, "Port number")

type echoHandler struct{}

func (e echoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n, err := io.Copy(w, r.Body)
	var errStr string
	if err != nil {
		errStr = " ERR:" + err.Error()
	}
	log.Printf("%s %s %s %d%s", r.RemoteAddr, r.Method, r.URL.Path, n, errStr)
}

func realMain() error {
	flag.Parse()

	sigChan := make(chan os.Signal, 10)
	signal.Notify(sigChan, os.Interrupt)

	addr := fmt.Sprintf("127.0.0.1:%d", *flagPort)
	log.Printf("Listening on %s", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)
	go func() {
		readBuf := make([]byte, 2048)
		for {
			c, err := l.Accept()
			if err != nil {
				errChan <- err
				return
			}

			log.Printf("Accepted connection from %s", c.RemoteAddr())
			for {
				n, err := c.Read(readBuf)
				if err != nil {
					log.Printf("TCP Read error: %v", err)
					break
				}

				_, err = c.Write(readBuf[:n])
				if err != nil {
					log.Printf("TCP Write error: %v", err)
					break
				}

				log.Printf("Echoed %d", n)
			}

		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-sigChan:
		log.Print("Received interrupt. Exiting.")
		return nil
	}
}

func main() {
	err := realMain()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
