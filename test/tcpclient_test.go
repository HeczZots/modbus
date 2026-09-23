// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.

package test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/HeczZots/modbus"
)

const (
	tcpDevice = "localhost:502"
)

func TestTCPClient(t *testing.T) {
	const slaveId = 1
	client := modbus.TCPClient(tcpDevice)
	ClientTestAll(t, client, slaveId)
}

func TestTCPClientAdvancedUsage(t *testing.T) {
	const slaveId = 1
	handler := modbus.NewTCPClientHandler(tcpDevice)
	handler.Timeout = 5 * time.Second
	handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
	if err := handler.Connect(); err != nil {
		t.Fatal(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadDiscreteInputs(slaveId, 15, 2)
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
	results, err = client.WriteMultipleRegisters(slaveId, 1, 2, []byte{0, 3, 0, 4})
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
	results, err = client.WriteMultipleCoils(slaveId, 5, 10, []byte{4, 3})
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
}

func TestTCPClientTwoSlaves(t *testing.T) {
	handler := modbus.NewTCPClientHandler(tcpDevice)
	handler.Timeout = 5 * time.Second
	handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
	if err := handler.Connect(); err != nil {
		t.Fatal(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	for _, slaveId := range []byte{1, 2} {
		results, err := client.ReadHoldingRegisters(slaveId, 0, 2)
		if err != nil {
			t.Fatalf("slave %d: %v", slaveId, err)
		}
		if len(results) != 4 {
			t.Fatalf("slave %d: expected 4 bytes, actual %d", slaveId, len(results))
		}
	}
}
