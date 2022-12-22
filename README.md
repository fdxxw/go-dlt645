# GO-DLT645
DLT645-2007 协议客户端 golang 实现

## 快速开始

```go
handler := NewDLT645ClientHandler("COM3")
handler.Baud = 2400 // 波特率
handler.StopBits = 1
handler.Parity = serial.ParityEven
handler.Size = 8
handler.ReadTimeout = 1 * time.Second
handler.Logger = log.New(os.Stdout, "rs485: ", log.LstdFlags) // 日志
client := NewClient(handler)
// 读取地址
addr, err := client.ReadAddr()
// 设置地址
err := client.SetAddr("202208310002")
// 读取数据
r, err := client.ReadData("202208310002", "028022FF")
// 设置数据
r, err := client.SetData("202208310002", "028022FF", 0x1C, ...)
// 发送 16 进制报文
r, err := client.SendHex("FEFEFEFE68AAAAAAAAAAAA681300DF16")

```
## License

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

