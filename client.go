package dlt645

// ClientHandler is the interface that groups the Packager and Transporter methods.
type ClientHandler interface {
	Packager
	Transporter
}

type client struct {
	packager    Packager
	transporter Transporter
}

// NewClient creates a new modbus client with given backend handler.
func NewClient(handler ClientHandler) Client {
	return &client{packager: handler, transporter: handler}
}

// NewClient2 creates a new modbus client with given backend packager and transporter.
func NewClient2(packager Packager, transporter Transporter) Client {
	return &client{packager: packager, transporter: transporter}
}

func (mb *client) ReadAddr() (results string, err error) {
	pdu := NewCommonProtocolDataUnit("AAAAAAAAAAAA", "13", "")
	res, err := mb.send(&pdu)
	if err != nil {
		return
	}
	data, err := res.Result(pdu.C)
	if err != nil {
		return
	}
	results = Byte2Hex(ByteRev(data))
	return
}

func (mb *client) SetAddr(addr string) (err error) {
	pdu := NewCommonProtocolDataUnit2(Hex2Byte("AAAAAAAAAAAA"), 0x15, ByteRev(Hex2Byte(addr)))
	res, err := mb.send(&pdu)
	if err != nil {
		return
	}
	_, err = res.Result(pdu.C)
	if err != nil {
		return
	}
	return
}
func (mb *client) ReadData(addr string, id string) (results []byte, err error) {
	pdu := NewCommonProtocolDataUnit2(ByteRev(Hex2Byte(addr)), 0x11, ByteRev(Hex2Byte(id)))
	res, err := mb.send(&pdu)
	if err != nil {
		return
	}
	results, err = res.Result(pdu.C)
	if err != nil {
		return
	}
	return
}
func (mb *client) SetData(addr string, id string, c byte, data []byte) (results []byte, err error) {
	pdu := NewCommonProtocolDataUnit2(ByteRev(Hex2Byte(addr)), c, append(ByteRev(Hex2Byte(id)), data...))
	res, err := mb.send(&pdu)
	if err != nil {
		return
	}
	results, err = res.Result(pdu.C)
	if err != nil {
		return
	}
	return
}

func (mb *client) SendHex(hex string) (results []byte, err error) {
	pdu := NewProtocolDataUnit(Hex2Byte(hex))
	res, err := mb.send(&pdu)
	if err != nil {
		return
	}
	results, err = res.Result(pdu.C)
	if err != nil {
		return
	}
	return
}
func (mb *client) Open() (err error) {
	err = mb.transporter.Open()
	return
}
func (mb *client) Close() (err error) {
	err = mb.transporter.Close()
	return
}
func (mb *client) send(request *ProtocolDataUnit) (response *ProtocolDataUnit, err error) {
	aduRequest, err := mb.packager.Encode(request)
	if err != nil {
		return
	}
	aduResponse, err := mb.transporter.Send(aduRequest)
	if err != nil {
		return
	}
	if err = mb.packager.Verify(aduRequest, aduResponse); err != nil {
		return
	}
	response, err = mb.packager.Decode(aduResponse)
	if err != nil {
		return
	}
	return
}
