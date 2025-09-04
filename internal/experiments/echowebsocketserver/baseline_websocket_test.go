// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"matheusd.com/testctx"
)

type goserbenchSmallStruct struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

// BenchmarkWsRPCCall benchmarks a baseline Websockets-based RPC system.
func BenchmarkWsRPCCall(b *testing.B) {
	sst := goserbenchSmallStruct{
		BirthDay: time.Now(),
		Siblings: 0x66669999,
		Spouse:   true,
		Money:    math.Float64frombits(0xabcd0000ef01),
		Name:     "slimshady0123456",
		Phone:    "phone67890",
	}

	sstAsJson, err := json.Marshal(sst)
	if err != nil {
		panic(err)
	}
	sstAsBytes := []byte(sstAsJson)

	b.Run("in-process server static data", func(b *testing.B) {
		// Run he server.
		errChan := make(chan error)
		var svr http.Server
		svr.Addr = fmt.Sprintf("127.0.0.1:8181")
		svr.Handler = echoHandler{skipLog: true}
		go func() {
			errChan <- svr.ListenAndServe()
		}()
		b.Cleanup(func() {
			svr.Shutdown(testctx.New(b))
			<-errChan
		})
		time.Sleep(10 * time.Millisecond)

		url := "ws://" + svr.Addr + "/echo"
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			b.Fatalf("Dial error: %v", err)
		}
		defer conn.Close()

		var lastMsg []byte

		b.SetBytes(int64(len(sstAsBytes)))
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			err := conn.WriteMessage(websocket.TextMessage, sstAsBytes)
			if err != nil {
				b.Fatal(err)
			}

			_, lastMsg, err = conn.ReadMessage()
			if err != nil {
				b.Fatal(err)
			}
		}

		if !bytes.Equal(lastMsg, sstAsBytes) {
			b.Fatal("wrong echo")
		}
	})

	b.Run("in-process server static data parallel", func(b *testing.B) {
		// Run he server.
		errChan := make(chan error)
		var svr http.Server
		svr.Addr = fmt.Sprintf("127.0.0.1:8181")
		svr.Handler = echoHandler{skipLog: true}
		go func() {
			errChan <- svr.ListenAndServe()
		}()
		b.Cleanup(func() {
			svr.Shutdown(testctx.New(b))
			<-errChan
		})
		time.Sleep(10 * time.Millisecond)

		url := "ws://" + svr.Addr + "/echo"

		b.SetBytes(int64(len(sstAsBytes)))
		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			conn, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				b.Fatalf("Dial error: %v", err)
			}
			defer conn.Close()

			var lastMsg []byte
			for pb.Next() {
				err := conn.WriteMessage(websocket.TextMessage, sstAsBytes)
				if err != nil {
					b.Fatal(err)
				}

				_, lastMsg, err = conn.ReadMessage()
				if err != nil {
					b.Fatal(err)
				}
			}

			if lastMsg != nil && !bytes.Equal(lastMsg, sstAsBytes) {
				b.Fatal("wrong echo")
			}
		})
	})

}
