package detector

import (
	"github.com/zhanglp0129/goproxypool/config"
	"github.com/zhanglp0129/goproxypool/storage"
	"net/http"
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

// 发请求检测连通性
func request(client http.Client, url string) error {
	_, err := client.Get(url)
	// TODO 记录 info 日志
	return err
}
