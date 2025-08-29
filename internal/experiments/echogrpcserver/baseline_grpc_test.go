// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	cmd := exec.Command("go", "run", ".", "-port", "8585")
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

	sst := &SmallStruct{
		BirthDay: time.Now().UnixNano(),
		Siblings: 0x66669999,
		Spouse:   true,
		Money:    math.Float64frombits(0xabcd0000ef01),
		Name:     "slimshady0123456",
		Phone:    "phone67890",
	}

	addr := fmt.Sprintf("127.0.0.1:8585")

	// Standard grpc API.
	b.Run("standard", func(b *testing.B) {
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
		conn, err := grpc.NewClient(addr, opts...)
		require.NoError(b, err)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		client := NewEchoClient(conn)
		_, err = client.Echo(ctx, sst)
		require.NoError(b, err)

		var res *SmallStruct

		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			res, err = client.Echo(ctx, sst)
			if err != nil {
				b.Fatal(err)
			}
		}

		require.Equal(b, sst.String(), res.String())
	})
}
