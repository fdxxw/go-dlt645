package dlt645

import (
	"io"
	"time"

	"github.com/spf13/cast"
)

type DLT645ClientHandler struct {
	dlt645SerialPackager
	dlt645SerialTransporter
}

// NewRTUClientHandler allocates and initializes a RTUClientHandler.
func NewDLT645ClientHandler(address string) *DLT645ClientHandler {
	handler := &DLT645ClientHandler{}
	handler.Name = address
	handler.ReadTimeout = serialTimeout
	handler.IdleTimeout = serialIdleTimeout
	return handler
}

// RTUClient creates RTU client with default handler and given connect string.
func DLT645Client(address string) Client {
	handler := NewDLT645ClientHandler(address)
	return NewClient(handler)
}

type dlt645SerialTransporter struct {
	serialPort
}
type dlt645SerialPackager struct{}

func (mb *dlt645SerialPackager) Encode(pdu *ProtocolDataUnit) (adu []byte, err error) {
	adu = pdu.Value()
	return
}

func (mb *dlt645SerialPackager) Decode(adu []byte) (pdu *ProtocolDataUnit, err error) {
	v := NewProtocolDataUnit(adu)
	pdu = &v
	return
}

func (mb *dlt645SerialPackager) Verify(aduRequest []byte, aduResponse []byte) (err error) {
	return
}

func (mb *dlt645SerialTransporter) Send(aduRequest []byte) (aduResponse []byte, err error) {
	// Make sure port is connected
	if err = mb.serialPort.connect(); err != nil {
		return
	}
	// Start the timer to close when idle
	mb.serialPort.lastActivity = time.Now()
	mb.serialPort.startCloseTimer()

	// 发送报文
	mb.serialPort.logf("dlt645: sending % x\n", aduRequest)
	if _, err = mb.port.Write(aduRequest); err != nil {
		return
	}
	// 延时在 20ms <= Td <= 500ms
	// 读取报文
	// 1. 先读取 14 个字节 	fefefefe 68 190002031122 68 91 84
	var data [1024]byte
	var n int
	var n1 int

	n, err = ReadAtLeast(mb.port, data[:], 14, 500*time.Millisecond)
	if err != nil {
		return
	}
	// 帧起始符长度
	frontLen := 0
	for _, b := range data {
		if b == 0xfe {
			frontLen++
		} else {
			break
		}
	}
	L := cast.ToInt(data[frontLen+9]) // 数据域长度
	// 总字节数
	bytesToRead := frontLen + 1 + 6 + 1 + 1 + 1 + L + 1 + 1
	// 读取剩余字节
	if n < bytesToRead {
		if bytesToRead > n {
			n1, err = io.ReadFull(mb.port, data[n:bytesToRead])
			n += n1
		}
	}
	aduResponse = data[:n]
	mb.serialPort.logf("dlt645: received % x\n", aduResponse)
	return
}

type D645Param struct {
	Identify string `json:"identify"` // 数据标识, 0280FF01
	Address  uint16 `json:"address"`  // 数据地址
	Quantity uint16 `json:"quantity"` // 字节个数
}

func (a *D645Param) Set(cfg map[string]interface{}) {
	a.Identify = cast.ToString(cfg["identify"])
	a.Address = cast.ToUint16(cfg["address"])
	a.Quantity = cast.ToUint16(cfg["quantity"])
}
