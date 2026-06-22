package wolconf

import (
	"fmt"
	"log"

	"github.com/gdy666/lucky/thirdlib/gdylib/netinterfaces"
	"github.com/gdy666/lucky/thirdlib/gdylib/stringsp"
	gowol "github.com/gdy666/lucky/thirdlib/go-wol"
)

var httpClientSecureVerify bool
var httpClientTimeout int

func GetHttpClientSecureVerify() bool {
	return httpClientSecureVerify
}

func SetHttpClientSecureVerify(b bool) {
	httpClientSecureVerify = b
}

func SetHttpClientTimeout(t int) {
	httpClientTimeout = t
}

type WOLDevice struct {
	Key          string
	DeviceName   string
	MacList      []string
	BroadcastIPs []string
	Port         int
	Relay        bool //交给中继设备发送
	Repeat       int  //重复发送次数
	//PowerOffCMD  string //关机指令
func (d *WOLDevice) GetIdentKey() string {
	return fmt.Sprintf("WOL:%s", d.Key)
}

func (d *WOLDevice) WakeUp(finishedCallback func(relay bool, macList []string, broadcastIps []string, port, repeat int)) error {
	return WakeOnLan(d.Relay, d.MacList, d.BroadcastIPs, d.Port, d.Repeat, finishedCallback)
}

func WakeOnLan(relay bool, macList []string, broadcastIps []string, port, repeat int,
	finishedCallback func(relay bool, macList []string, broadcastIps []string, port, repeat int),
) (err error) {
	defer func() {
		if finishedCallback != nil {
			finishedCallback(relay, macList, broadcastIps, port, repeat)
		}
	}()
	globalBroadcastList := netinterfaces.GetGlobalIPv4BroadcastList()
	matchCount := 0

	defer func() {
		if matchCount <= 0 {
			err = fmt.Errorf("找不到匹配的局域网广播IP,未能发送唤醒指令")
		}
	}()

	if len(broadcastIps) > 0 {
		for _, bcst := range broadcastIps {
			bcstOk := stringsp.StrIsInList(bcst, globalBroadcastList)
			if !bcstOk {
				continue
			}
			matchCount++
			for _, mac := range macList {
				gowol.WakeUpRepeat(mac, bcst, "", port, repeat)
			}

		}
		return
	}

	for _, bcst := range globalBroadcastList {
		matchCount++
		for _, mac := range macList {
			gowol.WakeUpRepeat(mac, bcst, "", port, repeat)
		}
	}

	return
}
