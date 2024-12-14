package detector

import (
	"errors"
	"github.com/zhanglp0129/goproxypool/common/pojo"
	"math/rand"
	"net/http"
	"time"
)

// 运行代理地址的可用性检测
func runAddressDetect() {
	// 开启一个计时器，每隔一段时间检测一次
	interval := time.Duration(CFG.Detect.Interval) * time.Second
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		// 获取指定个数的代理地址
		addresses, err := Storage.GetDetectedProxyAddresses()
		if err != nil {
			// TODO 记录日志
			continue
		}
		// 开启goroutine，执行检测
		for _, address := range addresses {
			go func() {
				err = detect(address)
				if err != nil {
					// TODO 记录 info 日志
				}
				// 检测完成
				err = Storage.FinishDetection(address.ID, err == nil)
				if err != nil {
					// TODO 记录日志
				}
			}()
		}
	}
}

// 检测代理地址的可用性
func detect(address pojo.ProxyAddress) error {
	// 获取重试次数
	attempts := CFG.Detect.Attempts
	// 返回的error
	var res error
	for i := 0; i < attempts; i++ {
		var website string
		func() {
			websitesMutex.RLock()
			defer websitesMutex.RUnlock()
			website = availableWebsites[rand.Intn(len(availableWebsites))]
		}()
		// 向网站发送请求，检测代理地址可用性
		client := http.Client{
			Timeout: time.Duration(CFG.Detect.Timeout) * time.Second,
		}
		err := request(client, website)
		if err != nil {
			res = errors.Join(res, err)
		} else {
			return nil
		}
	}
	return res
}
