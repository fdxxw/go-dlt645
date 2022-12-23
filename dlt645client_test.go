package dlt645

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/fdxxw/go-wen"
	"github.com/tarm/serial"
)

func TestSend(t *testing.T) {
	transporter := &dlt645SerialTransporter{}
	transporter.Baud = 2400
	transporter.StopBits = 1
	transporter.Parity = serial.ParityEven
	transporter.Size = 8
	transporter.ReadTimeout = 1 * time.Second
	transporter.Name = "COM3"
	frame := ProtocolDataUnit{
		Front:   []byte{0xfe, 0xfe, 0xfe, 0xfe},
		Start:   0x68,
		Address: Hex2Byte("AAAAAAAAAAAA"),
		C:       Hex2Byte("13")[0],
		Data:    []byte{},
		End:     0x16,
	}
	defer wen.TimeCost()()
	for i := 0; i < 1; i++ {
		_, err := transporter.Send(frame.Value())
		if err != nil {
			t.Fatal(err)
		}
		// log.Println(r)
	}
}

func newTestClient() Client {
	handler := NewDLT645ClientHandler("COM3")
	handler.Baud = 2400
	handler.StopBits = 1
	handler.Parity = serial.ParityEven
	handler.Size = 8
	handler.ReadTimeout = 1 * time.Second
	handler.Logger = log.New(os.Stdout, "rs485: ", log.LstdFlags)
	client := NewClient(handler)
	return client
}

func TestReadAddr(t *testing.T) {
	client := newTestClient()
	addr, err := client.ReadAddr()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(addr)
}
func TestSetAddr(t *testing.T) {
	client := newTestClient()
	err := client.SetAddr("202208310002")
	if err != nil {
		t.Fatal(err)
	}
}
func TestReadData(t *testing.T) {
	client := newTestClient()
	r, err := client.ReadData("202208310002", "028022FF")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(Byte2Hex(r))
}
func TestSendHex(t *testing.T) {
	client := newTestClient()
	r, err := client.SendHex("fe fe fe fe 68 02 00 31 08 22 20 68 11 04 32 55 b3 35 d1 16")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(Byte2Hex(ByteRev(r)))
}
