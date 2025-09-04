// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

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

func TestMain(m *testing.M) {
	cmd := exec.Command("go", "run", ".", "-port", "8080")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	time.Sleep(time.Second)
	if err != nil {
		panic(err)
	}
	m.Run()
	cmd.Process.Kill()
}

// BenchmarkHTTPRPCCall benchmarks a baseline HTTP RPC system.
func BenchmarkHTTPRPCCall(b *testing.B) {
	url := "http://127.0.0.1:8080/echo"

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

	// Best case scenario for data that is already somehow encoded.
	b.Run("default client static data", func(b *testing.B) {
		reqBody := bytes.NewReader(sstAsBytes)
		reply := make([]byte, len(sstAsBytes))

		b.SetBytes(int64(len(sstAsBytes)))
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			res, err := http.Post(url, "text/json", reqBody)
			if err != nil {
				b.Fatal(err)
			}
			reqBody.Reset(sstAsBytes)

			n, err := res.Body.Read(reply)
			if err != nil && !errors.Is(err, io.EOF) {
				b.Fatal(err)
			}
			if n != len(sstAsBytes) {
				b.Fatal("wrong number of bytes")
			}
		}

		if !bytes.Equal(reply, sstAsBytes) {
			b.Fatal("wrong echo")
		}
	})

	// Assuming the client needs to encode every message as json on every
	// call but that we can reuse the buffer somehow.
	b.Run("default client json encode", func(b *testing.B) {
		reqBuf := bytes.NewBuffer(make([]byte, 0, len(sstAsBytes)*2))
		enc := json.NewEncoder(reqBuf)
		reply := make([]byte, len(sstAsBytes))

		b.SetBytes(int64(len(sstAsBytes)))
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			if err := enc.Encode(sst); err != nil {
				b.Fatal(err)
			}

			res, err := http.Post(url, "text/json", reqBuf)
			if err != nil {
				b.Fatal(err)
			}
			reqBuf.Reset()

			n, err := res.Body.Read(reply)
			if err != nil && !errors.Is(err, io.EOF) {
				b.Fatal(err)
			}
			if n != len(sstAsBytes) {
				b.Fatal("wrong number of bytes")
			}
		}

		if !bytes.Equal(reply, sstAsBytes) {
			b.Fatal("wrong echo")
		}
	})

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

		reqBody := bytes.NewReader(sstAsBytes)
		reply := make([]byte, len(sstAsBytes))

		url := "http://" + svr.Addr + "/echo"

		b.SetBytes(int64(len(sstAsBytes)))

		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			res, err := http.Post(url, "text/json", reqBody)
			if err != nil {
				b.Fatal(err)
			}
			reqBody.Reset(sstAsBytes)

			n, err := res.Body.Read(reply)
			if err != nil && !errors.Is(err, io.EOF) {
				b.Fatal(err)
			}
			if n != len(sstAsBytes) {
				b.Fatal("wrong number of bytes")
			}
		}

		if !bytes.Equal(reply, sstAsBytes) {
			b.Fatal("wrong echo")
		}

	})

	b.Run("in-process server static data parallel", func(b *testing.B) {
		// Run he server.

		errChan := make(chan error)
		var svr http.Server
		svr.Addr = fmt.Sprintf("127.0.0.1:8282")
		svr.Handler = echoHandler{skipLog: true}
		go func() {
			errChan <- svr.ListenAndServe()
		}()
		b.Cleanup(func() {
			svr.Shutdown(testctx.New(b))
			<-errChan
		})
		time.Sleep(10 * time.Millisecond)

		url := "http://" + svr.Addr + "/echo"

		b.SetBytes(int64(len(sstAsBytes)))

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			gotOne := false

			var c http.Client
			dialer := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			c.Transport = &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				DialContext:           dialer.DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			}

			reqBody := bytes.NewReader(sstAsBytes)
			reply := make([]byte, len(sstAsBytes))
			for pb.Next() {
				res, err := c.Post(url, "text/json", reqBody)
				if err != nil {
					b.Fatal(err)
				}
				reqBody.Reset(sstAsBytes)

				n, err := res.Body.Read(reply)
				if err != nil && !errors.Is(err, io.EOF) {
					b.Fatal(err)
				}
				if n != len(sstAsBytes) {
					b.Fatal("wrong number of bytes")
				}
				gotOne = true
			}

			if gotOne && !bytes.Equal(reply, sstAsBytes) {
				b.Fatal("wrong echo")
			}
		})

	})
}
