package swctrl

import (
	"net"

	"github.com/astaxie/beego/logs"
)

var (
	SW_OPEN  byte = 0x0f //开
	SW_CLOSE byte = 0x10 //关
)

type SWCtrl struct {
	Conn net.Conn
}

func NewSWCtrl(adapter, urlAddr string) (pCtrl *SWCtrl, err error) {
	conn, err := net.Dial(adapter, urlAddr)
	pCtrl = nil

	if err == nil {
		pCtrl = &SWCtrl{Conn: conn}
	} else {
		logs.Error(err)
	}

	return
}

func (p *SWCtrl) Close() {
	p.Conn.Close()
}

// RecData 接受状态码
func (p *SWCtrl) RecData() ([]byte, error) {
	var rb = make([]byte, 20)
	_, err := p.Conn.Read(rb)

	if err != nil {
		logs.Error("返回错误:", err)
	} else {
		logs.Info("返回:", rb)
	}

	return rb, err
}

// SendData 发送控制码
func (p *SWCtrl) SendData(sb []byte) error {
	_, err := p.Conn.Write(sb)

	logs.Info("发送:", sb)

	if err != nil {
		logs.Error(err.Error())
	}

	return err
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

func init() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/swctrl/swctrl.log","level":7,"maxlines":0,"maxsize":0,"maxdays":10,"perm":"777"}`)
}
