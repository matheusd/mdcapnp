// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

type echoHandler struct {
	skipLog bool
}

func (e echoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if !e.skipLog {
			log.Println("Upgrade error:", err)
		}
		return
	}
	defer conn.Close()

	for {
		// Read message from client
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			if !e.skipLog {
				log.Println("Read error:", err)
			}
			return
		}

		// Echo message back to client
		err = conn.WriteMessage(msgType, msg)
		if err != nil {
			if !e.skipLog {
				log.Println("Write error:", err)
			}
			return
		}
	}

}
