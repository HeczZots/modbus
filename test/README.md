System testing for [modbus library](https://github.com/HeczZots/modbus)

Modbus simulator
----------------
*   [Diagslave](http://www.modbusdriver.com/diagslave.html)
*   [socat](http://www.dest-unreach.org/socat/)

```bash
# TCP
$ diagslave -m tcp -p 502

# RTU/ASCII
$ socat -d -d pty,raw,echo=0 pty,raw,echo=0
2015/04/03 12:34:56 socat[2342] N PTY is /dev/pts/6
2015/04/03 12:34:56 socat[2342] N PTY is /dev/pts/7
$ diagslave -m ascii /dev/pts/7

# Or
$ diagslave -m rtu /dev/pts/7

$ go test -v -run TCP
$ go test -v -run RTU
$ go test -v -run ASCII
```

`TestTCPClientTwoSlaves` expects a TCP server on port 502 that answers
slave ids 1 and 2. Port 502 may require administrator privileges.
