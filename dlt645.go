package dlt645

import "errors"

type ProtocolDataUnit struct {
	Front   []byte // 在主站发送帧信息之前，先发送1—4个字节FEH，以唤醒接收方。
	Start   byte   // 标识一帧信息的开始，其值为68H=01101000B。
	Address []byte // 地址域 地址域由 6 个字节构成，每字节 2 位 BCD 码 地址域传输时低字节在前，高字节在后。
	C       byte   // 控制码 C
	L       uint8  // 数据域长度  L为数据域的字节数。读数据时L≤200，写数据时L≤50，L=0表示无数据域
	Data    []byte // 数据域 数据域包括数据标识、密码、操作者代码、数据、帧序号等，其结构随控制码的功能而改变。传输时发送方按字节进行加33H处理，接收方按字节进行减33H处理。
	CS      byte   // 校验码 从第一个帧起始符开始到校验码之前的所有各字节的模256的和，即各字节二进制算术和，不计超过256的溢出值
	End     byte   // 标识一帧信息的结束，其值为16H=00010110B。
}

// 16 进制字符串 帧
func (a *ProtocolDataUnit) Value() []byte {
	a.L = uint8(len(a.Data))
	// 计算 CS
	bs := []byte{0x68}              // 起始符
	bs = append(bs, a.Address...)   // 地址域
	bs = append(bs, 0x68, a.C, a.L) // 控制码和数据域长度
	bs = append(bs, a.Data...)      // 数据域
	a.CS = a.computeCS(bs)
	finalBs := []byte{}
	finalBs = append(finalBs, a.Front...)
	finalBs = append(finalBs, bs...)
	finalBs = append(finalBs, a.CS, a.End)
	return finalBs
}

func (a *ProtocolDataUnit) StringValue() string {
	return Byte2Hex(a.Value())
}

func (a *ProtocolDataUnit) computeCS(data []byte) byte {
	var sum uint64
	for _, b := range data {
		sum += uint64(b)
	}
	return byte(sum % 256)
}

// 数据标识
func (a *ProtocolDataUnit) Identify() string {
	if a.L >= 4 {
		return Byte2Hex(ByteRev(a.Data[:4]))
	} else {
		return ""
	}
}

// 根据前一个控制码判断成功失败，返回数据（减去 33H 之后的）
func (a *ProtocolDataUnit) Result(cmdC byte) (data []byte, err error) {
	errC := cmdC + 0xC0
	if a.C == errC {
		err = errors.New("error")
	}
	data = a.Data
	return
}

func NewCommonProtocolDataUnit(addr string, c string, data string) ProtocolDataUnit {
	return ProtocolDataUnit{
		Front:   []byte{0xfe, 0xfe, 0xfe, 0xfe},
		Start:   0x68,
		End:     0x16,
		Address: Hex2Byte(addr),
		C:       Hex2Byte(c)[0],
		Data:    ByteAdd(Hex2Byte(data), 0x33),
	}
}
func NewCommonProtocolDataUnit2(addr []byte, c byte, data []byte) ProtocolDataUnit {
	return ProtocolDataUnit{
		Front:   []byte{0xfe, 0xfe, 0xfe, 0xfe},
		Start:   0x68,
		End:     0x16,
		Address: addr,
		C:       c,
		Data:    ByteAdd(data, 0x33),
	}
}
func NewProtocolDataUnit(bs []byte) ProtocolDataUnit {
	frame := ProtocolDataUnit{
		Start: 0x68,
		End:   0x16,
	}
	for _, b := range bs {
		if b == 0xfe {
			frame.Front = append(frame.Front, b)
		} else {
			break
		}
	}
	bs = bs[len(frame.Front):] // 删除前导向
	// 帧起始符 68H
	if bs[0] == 0x68 {
		bs = bs[1:]
		// 地址域
		frame.Address = ByteRev(bs[:6])
		bs = bs[6:]
		bs = bs[1:]     // 帧起始符
		frame.C = bs[0] // 控制码
		frame.L = bs[1] // 数据域长度
		bs = bs[2:]
		data := bs[:frame.L] // 数据域
		for _, d := range data {
			frame.Data = append(frame.Data, d-0x33) // 减 33H
		}
		bs = bs[frame.L:]
		frame.CS = bs[0]  // 校验码
		frame.End = bs[1] // 帧结束符
	}
	return frame
}

type Packager interface {
	Encode(pdu *ProtocolDataUnit) (adu []byte, err error)
	Decode(adu []byte) (pdu *ProtocolDataUnit, err error)
	Verify(aduRequest []byte, aduResponse []byte) (err error)
}

// Transporter specifies the transport layer.
type Transporter interface {
	Send(aduRequest []byte) (aduResponse []byte, err error)
	Open() (err error)
	Close() (err error)
}
