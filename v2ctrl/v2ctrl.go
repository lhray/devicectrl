package v2ctrl

import (
	"fmt"
	"net"

	"github.com/astaxie/beego/logs"
)

var (
	V2_OPEN  byte = 1 //开
	V2_CLOSE byte = 0 //关
)

type V2Ctrl struct {
	Conn net.Conn
}

func NewV2Ctrl(adapter, urlAddr string) (pCtrl *V2Ctrl, err error) {
	conn, err := net.Dial(adapter, urlAddr)
	pCtrl = nil

	if err == nil {
		pCtrl = &V2Ctrl{Conn: conn}
	} else {
		logs.Error(err)
	}

	return
}

func (p *V2Ctrl) Close() {
	p.Conn.Close()
}

// RecData 接受状态码
func (p *V2Ctrl) RecData() ([]byte, error) {
	var rb = make([]byte, 1024)
	_, err := p.Conn.Read(rb)

	if err != nil {
		logs.Error("返回错误:", err)
	} else {
		logs.Info("返回:", string(rb))
	}

	return rb, err
}

// SendData 发送控制码
func (p *V2Ctrl) SendData(sb []byte) error {
	_, err := p.Conn.Write(sb)

	logs.Info("发送:", string(sb))

	if err != nil {
		logs.Error(err.Error())
	}

	return err
}

// MakeCodes 构造控制码
// actionType: SW_OPEN, SW_CLOSE
func (p *V2Ctrl) MakeCodes(channel int, actionType byte) []byte {
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

func (p *V2Ctrl) OpenWithClose(channel, relaySecond int) []byte {
	str := fmt.Sprintf("on%d:%02d", channel, relaySecond)

	return []byte(str)
}

func (p *V2Ctrl) GetStatus(channel int) string {
	str := fmt.Sprintf("read%d", channel)

	p.SendData([]byte(str))
	sb, _ := p.RecData()

	return string(sb)
}

func init() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/v2ctrl/v2ctrl.log","level":7,"maxlines":0,"maxsize":0,"maxdays":10,"perm":"777"}`)
}
