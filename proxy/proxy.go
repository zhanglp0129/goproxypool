package proxy

import (
	"github.com/zhanglp0129/goproxypool/config"
	"github.com/zhanglp0129/goproxypool/storage"
)

var (
	CFG     = config.CFG
	Storage = storage.Storage
)

// Run 在后台运行代理服务器
func Run() {
	go runHttp()
}
