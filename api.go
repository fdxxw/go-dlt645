package dlt645

type Client interface {
	// 读表地址
	ReadAddr() (results string, err error)
	// 设置表地址
	SetAddr(addr string) error

	// 读数据
	// addr 地址
	// id 数据标识符
	ReadData(addr string, id string) (results []byte, err error)

	// 设置
	// addr 地址
	// id 数据标识符
	// c 控制码
	// data 数据，未加 33H 的
	SetData(addr string, id string, c byte, data []byte) (results []byte, err error)

	// 发送 16进制字符串，发送原始报文
	SendHex(hex string) (results []byte, err error)
	// 打开端口
	Open() (err error)
	// 关闭
	Close() (err error)
}
