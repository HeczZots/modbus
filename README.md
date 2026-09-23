go modbus [![GoDoc](https://pkg.go.dev/badge/github.com/HeczZots/modbus.svg)](https://pkg.go.dev/github.com/HeczZots/modbus)
=========
Fault-tolerant, fail-fast implementation of Modbus protocol in Go.

Supported functions
-------------------
Bit access:
*   Read Discrete Inputs
*   Read Coils
*   Write Single Coil
*   Write Multiple Coils

16-bit access:
*   Read Input Registers
*   Read Holding Registers
*   Write Single Register
*   Write Multiple Registers
*   Read/Write Multiple Registers
*   Mask Write Register
*   Read FIFO Queue

Supported formats
-----------------
*   TCP
*   Serial (RTU, ASCII)

Usage
-----
Install:
```bash
go get github.com/HeczZots/modbus
```

Every request takes the slave id (unit identifier) as its first argument,
so one handler can talk to several devices on the same bus.

Basic usage:
```go
// Modbus TCP
client := modbus.TCPClient("localhost:502")
// Read input register 9 from slave 1
results, err := client.ReadInputRegisters(1, 8, 1)

// Modbus RTU/ASCII
// Default configuration is 19200, 8, 1, even
client = modbus.RTUClient("/dev/ttyS0")
results, err = client.ReadCoils(1, 2, 1)
```

Advanced usage:
```go
// Modbus TCP
handler := modbus.NewTCPClientHandler("localhost:502")
handler.Timeout = 10 * time.Second
handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
// Connect manually so that multiple requests are handled in one connection session
err := handler.Connect()
defer handler.Close()

client := modbus.NewClient(handler)
results, err := client.ReadDiscreteInputs(1, 15, 2)
results, err = client.WriteMultipleRegisters(1, 1, 2, []byte{0, 3, 0, 4})
results, err = client.WriteMultipleCoils(1, 5, 10, []byte{4, 3})

// Several devices over one connection
results, err = client.ReadHoldingRegisters(1, 0, 2)
results, err = client.ReadHoldingRegisters(2, 0, 2)
```

```go
// Modbus RTU/ASCII
handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
handler.BaudRate = 115200
handler.DataBits = 8
handler.Parity = "N"
handler.StopBits = 1
handler.Timeout = 5 * time.Second

err := handler.Connect()
defer handler.Close()

client := modbus.NewClient(handler)
results, err := client.ReadDiscreteInputs(1, 15, 2)
```

Breaking change
---------------
Previous versions stored the slave id in the handler (`handler.SlaveId`).
That field is gone: every `Client` method now takes the slave id as its
first argument, so custom `Client` implementations and mocks must be
updated. Custom `Packager` implementations must update `Encode` to
`Encode(slaveId byte, pdu *ProtocolDataUnit) (adu []byte, err error)`.

References
----------
-   [Modbus Specifications and Implementation Guides](http://www.modbus.org/specs.php)
