// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package modbus

import (
	"bytes"
	"testing"
)

type fakeHandler struct {
	rtuPackager
	slaveIds []byte
	response []byte
}

func (h *fakeHandler) Encode(slaveId byte, pdu *ProtocolDataUnit) ([]byte, error) {
	h.slaveIds = append(h.slaveIds, slaveId)
	return h.rtuPackager.Encode(slaveId, pdu)
}

func (h *fakeHandler) Send(aduRequest []byte) ([]byte, error) {
	return h.response, nil
}

func TestClientPassesSlaveIdToPackager(t *testing.T) {
	pdu := ProtocolDataUnit{FunctionCode: FuncCodeReadHoldingRegisters, Data: []byte{4, 0, 0x0A, 1, 2}}
	var p rtuPackager
	response, err := p.Encode(0x11, &pdu)
	if err != nil {
		t.Fatal(err)
	}
	handler := &fakeHandler{response: response}
	client := NewClient(handler)

	results, err := client.ReadHoldingRegisters(0x11, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	if expected := []byte{0, 0x0A, 1, 2}; !bytes.Equal(expected, results) {
		t.Fatalf("results: expected %v, actual %v", expected, results)
	}
	if expected := []byte{0x11}; !bytes.Equal(expected, handler.slaveIds) {
		t.Fatalf("slave ids passed to packager: expected %v, actual %v", expected, handler.slaveIds)
	}

	_, err = client.ReadHoldingRegisters(0x22, 0, 2)
	if err == nil {
		t.Fatal("expected slave id mismatch error")
	}
	if expected := []byte{0x11, 0x22}; !bytes.Equal(expected, handler.slaveIds) {
		t.Fatalf("slave ids passed to packager: expected %v, actual %v", expected, handler.slaveIds)
	}
}

func TestClientSlaveIdPerFunction(t *testing.T) {
	calls := []struct {
		name    string
		slaveId byte
		call    func(c Client, slaveId byte) error
	}{
		{"ReadCoils", 1, func(c Client, id byte) error { _, err := c.ReadCoils(id, 0, 1); return err }},
		{"ReadDiscreteInputs", 2, func(c Client, id byte) error { _, err := c.ReadDiscreteInputs(id, 0, 1); return err }},
		{"WriteSingleCoil", 3, func(c Client, id byte) error { _, err := c.WriteSingleCoil(id, 0, 0xFF00); return err }},
		{"WriteMultipleCoils", 4, func(c Client, id byte) error { _, err := c.WriteMultipleCoils(id, 0, 1, []byte{1}); return err }},
		{"ReadInputRegisters", 5, func(c Client, id byte) error { _, err := c.ReadInputRegisters(id, 0, 1); return err }},
		{"ReadHoldingRegisters", 6, func(c Client, id byte) error { _, err := c.ReadHoldingRegisters(id, 0, 1); return err }},
		{"WriteSingleRegister", 7, func(c Client, id byte) error { _, err := c.WriteSingleRegister(id, 0, 1); return err }},
		{"WriteMultipleRegisters", 8, func(c Client, id byte) error { _, err := c.WriteMultipleRegisters(id, 0, 1, []byte{0, 1}); return err }},
		{"ReadWriteMultipleRegisters", 9, func(c Client, id byte) error {
			_, err := c.ReadWriteMultipleRegisters(id, 0, 1, 0, 1, []byte{0, 1})
			return err
		}},
		{"MaskWriteRegister", 10, func(c Client, id byte) error { _, err := c.MaskWriteRegister(id, 0, 1, 2); return err }},
		{"ReadFIFOQueue", 11, func(c Client, id byte) error { _, err := c.ReadFIFOQueue(id, 0); return err }},
	}
	for _, c := range calls {
		t.Run(c.name, func(t *testing.T) {
			handler := &fakeHandler{}
			client := NewClient(handler)
			if err := c.call(client, c.slaveId); err == nil {
				t.Fatalf("%s: expected error on empty response", c.name)
			}
			if expected := []byte{c.slaveId}; !bytes.Equal(expected, handler.slaveIds) {
				t.Fatalf("slave ids passed to packager: expected %v, actual %v", expected, handler.slaveIds)
			}
		})
	}
}
