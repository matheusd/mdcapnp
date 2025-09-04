// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"matheusd.com/depvendoredtestify/require"
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
	cmd := exec.Command("go", "run", ".", "-port", "9090")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	time.Sleep(time.Second)
	if err != nil {
		panic(err)
	}
	m.Run()
	cmd.Process.Kill()
}

// BenchmarkTCPRPCCall benchmarks a baseline TCP RPC system.
func BenchmarkTCPRPCCall(b *testing.B) {

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

	addr := fmt.Sprintf("127.0.0.1:9090")

	// Best case scenario for data that is already somehow encoded.
	b.Run("static data", func(b *testing.B) {
		reply := make([]byte, len(sstAsBytes))

		conn, err := net.Dial("tcp", addr)
		require.NoError(b, err)
		defer conn.Close()

		b.SetBytes(int64(len(sstAsBytes)))
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			_, err := conn.Write(sstAsBytes)
			if err != nil {
				b.Fatal(err)
			}

			_, err = io.ReadFull(conn, reply)
			if err != nil {
				b.Fatal(err)
			}
		}

		if !bytes.Equal(reply, sstAsBytes) {
			b.Fatal("wrong echo")
		}
	})

	// Assuming the client needs to encode every message as json on every
	// call but that we can reuse the buffer somehow.
	b.Run("json encode", func(b *testing.B) {
		conn, err := net.Dial("tcp", addr)
		require.NoError(b, err)

		reply := make([]byte, len(sstAsBytes))
		reqBuf := bytes.NewBuffer(make([]byte, 0, len(sstAsBytes)*2))
		enc := json.NewEncoder(reqBuf)

		b.SetBytes(int64(len(sstAsBytes)))
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			if err := enc.Encode(sst); err != nil {
				b.Fatal(err)
			}

			_, err := conn.Write(sstAsBytes)
			if err != nil {
				b.Fatal(err)
			}

			_, err = io.ReadFull(conn, reply)
			if err != nil {
				b.Fatal(err)
			}

		}

		if !bytes.Equal(reply, sstAsBytes) {
			b.Fatal("wrong echo")
		}
	})

	b.Run("in-process server static data", func(b *testing.B) {
		addr := "127.0.0.1:9191"
		l, err := net.Listen("tcp", addr)
		if err != nil {
			b.Fatal(err)
		}
		b.Cleanup(func() { l.Close() })
		svr := echoTcpServer{l: l, skipLog: true}
		go svr.run()

		reply := make([]byte, len(sstAsBytes))

		conn, err := net.Dial("tcp", addr)
		require.NoError(b, err)
		defer conn.Close()

		b.SetBytes(int64(len(sstAsBytes)))
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			_, err := conn.Write(sstAsBytes)
			if err != nil {
				b.Fatal(err)
			}

			_, err = io.ReadFull(conn, reply)
			if err != nil {
				b.Fatal(err)
			}
		}

		if !bytes.Equal(reply, sstAsBytes) {
			b.Fatal("wrong echo")
		}
	})

	b.Run("in-process server static data parallel bench", func(b *testing.B) {
		addr := "127.0.0.1:9191"
		l, err := net.Listen("tcp", addr)
		if err != nil {
			b.Fatal(err)
		}
		b.Cleanup(func() { l.Close() })
		svr := echoTcpServer{l: l, skipLog: true}
		go svr.run()

		b.SetBytes(int64(len(sstAsBytes)))
		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			reply := make([]byte, len(sstAsBytes))

			conn, err := net.Dial("tcp", addr)
			require.NoError(b, err)
			defer conn.Close()
			gotOne := false

			for pb.Next() {
				_, err := conn.Write(sstAsBytes)
				if err != nil {
					b.Fatal(err)
				}

				_, err = io.ReadFull(conn, reply)
				if err != nil {
					b.Fatal(err)
				}
				gotOne = true
			}

			if gotOne && !bytes.Equal(reply, sstAsBytes) {
				b.Fatalf("wrong echo: %q vs %q", reply, sstAsBytes)
			}
		})

	})

}
