// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package main

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative structdef.proto

import (
	context "context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

var flagPort = flag.Int("port", 0, "Port number")

type echoHandler struct {
	UnimplementedEchoServer
}

func (e echoHandler) Echo(_ context.Context, in *SmallStruct) (*SmallStruct, error) {
	log.Printf("Echoing back")
	return in, nil
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

	server := grpc.NewServer()
	RegisterEchoServer(server, echoHandler{})

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Serve(l)
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
