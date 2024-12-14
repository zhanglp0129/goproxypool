package detector

import (
	"github.com/zhanglp0129/goproxypool/config"
	"github.com/zhanglp0129/goproxypool/storage"
)

var (
	CFG     = config.CFG
	Storage = storage.Storage
)

// Run 在后台运行可用性检测
func Run() {
	go runAddressDetect()
	go runWebsiteDetect()
}
