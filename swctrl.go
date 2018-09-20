package devicectrl

import (
	"net"

	"github.com/astaxie/beego/logs"
)

var (
	SW_OPEN  byte = 0x0f //开
	SW_CLOSE byte = 0x10 //关
)

type SWCtrl struct{}

func NewSWCtrl() *SWCtrl {
	return new(SWCtrl)
}

// RecData 接受状态码
func (sw *SWCtrl) RecData(connection *net.TCPConn) []byte {
	var rb = make([]byte, 20)
	_, err := connection.Read(rb)

	if err != nil {
		logs.Error("返回错误:", err)
	} else {
		logs.Info("返回:", rb)
	}

	return rb
}

// SendData 发送控制码
func (sw *SWCtrl) SendData(conn *net.TCPConn, sb []byte) {
	_, err := conn.Write(sb)

	logs.Info("发送:", sb)

	if err != nil {
		logs.Error(err)
	}
}

// MakeCodes 构造控制码
// actionType: SW_OPEN, SW_CLOSE
func (sw *SWCtrl) MakeCodes(channel int, actionType byte) []byte {
	sb := []byte{0xaa, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xbb}

	if channel == 0 {
		// 全开全关
		if actionType == SW_OPEN {
			sb[1] = 0x0a //全开
		} else {
			sb[1] = 0x0b //全关
		}

	} else {
		// 单通道分别控制
		sb[1] = actionType // SW_OPEN or SW_CLOSE
		// 按通道控制
		sb[2] = byte(channel - 1) //转换为从0开始

		if actionType == SW_OPEN {
			sb[3] = 0x01 //开
		} else {
			sb[3] = 0x02 //关
		}
	}

	return sb
}
