package v2ctrl

import (
	"net"

	"github.com/astaxie/beego/logs"
)

var (
	V2_OPEN  byte = 1 //开
	V2_CLOSE byte = 0 //关
)

type V2Ctrl struct{}

func NewV2Ctrl() *V2Ctrl {
	return new(V2Ctrl)
}

// RecData 接受状态码
func (sw *V2Ctrl) RecData(connection *net.TCPConn) ([]byte, error) {
	var rb = make([]byte, 1024)
	_, err := connection.Read(rb)

	if err != nil {
		logs.Error("返回错误:", err)
	} else {
		logs.Info("返回:", string(rb))
	}

	return rb, err
}

// SendData 发送控制码
func (sw *V2Ctrl) SendData(conn *net.TCPConn, sb []byte) error {
	_, err := conn.Write(sb)

	logs.Info("发送:", string(sb))

	if err != nil {
		logs.Error(err.Error())
	}

	return err
}

// MakeCodes 构造控制码
// actionType: SW_OPEN, SW_CLOSE
func (sw *V2Ctrl) MakeCodes(channel int, actionType byte) []byte {
	var str string

	if channel == 1 {

		if actionType == V2_OPEN {
			str = "on1"
		} else if actionType == V2_CLOSE {
			str = "off1"
		}

	} else if channel == 2 {
		if actionType == V2_OPEN {
			str = "on2"
		} else if actionType == V2_CLOSE {
			str = "off2"
		}
	}

	return []byte(str)
}

func init() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/v2ctrl/v2ctrl.log","level":7,"maxlines":0,"maxsize":0,"maxdays":10,"perm":"777"}`)
}
