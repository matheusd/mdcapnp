// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var flagPort = flag.Int("port", 0, "Port number")

type echoHandler struct {
	skipLog bool
}

func (e echoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n, err := io.Copy(w, r.Body)
	var errStr string
	if err != nil {
		errStr = " ERR:" + err.Error()
	}
	if !e.skipLog {
		log.Printf("%s %s %s %d%s", r.RemoteAddr, r.Method, r.URL.Path, n, errStr)
	}
}

func realMain() error {
	flag.Parse()

	sigChan := make(chan os.Signal, 10)
	signal.Notify(sigChan, os.Interrupt)

	addr := fmt.Sprintf("127.0.0.1:%d", *flagPort)
	log.Printf("Listening on %s", addr)
	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(addr, echoHandler{})
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
