package detector

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	availableWebsites = make([]string, 0)
	websitesMutex     sync.RWMutex
)

// 运行直连网站连通性检测
func runWebsiteDetect() {
	// 开启计时器
	directInterval := time.Duration(CFG.Detect.DirectInterval) * time.Second
	ticker := time.NewTicker(directInterval)
	for {
		<-ticker.C
		doWebsiteDetect()
	}
}

// 执行直连网站检测
func doWebsiteDetect() {
	// 直连检测
	websites := CFG.Detect.Websites
	var mtx sync.Mutex
	tempWebsites := make([]string, 0)
	var wait sync.WaitGroup
	wait.Add(len(websites))
	for _, website := range websites {
		go func() {
			err := websiteDetect(website)
			if err != nil {
				// TODO 记录 warn 日志
				fmt.Printf("warn: 直连检测 %s 出错 %v\n", website, err)
			} else {
				// 检测成功，往临时网站切片中写入当前网站
				mtx.Lock()
				defer mtx.Unlock()
				tempWebsites = append(tempWebsites, website)
			}
			wait.Done()
		}()
	}
	wait.Wait()
	// 直连检测完成
	fmt.Printf("info: 直连检测完成，可用的检测网站 %v", tempWebsites)
	// 将临时网站的数据写入可用的网站切片
	func() {
		websitesMutex.Lock()
		defer websitesMutex.Unlock()
		availableWebsites = availableWebsites[0:0]
		for _, website := range tempWebsites {
			availableWebsites = append(availableWebsites, website)
		}
	}()
}

// 直连网站连通性检测
func websiteDetect(website string) error {
	// 获取重试次数
	attempts := CFG.Detect.Attempts
	var res error
	for i := 0; i < attempts; i++ {
		// 发请求判断连通性
		client := http.Client{
			Timeout: time.Duration(CFG.Detect.Timeout) * time.Second,
		}
		resp, err := client.Get(website)
		if err != nil {
			res = errors.Join(res, err)
		} else {
			// TODO 记录 info 日志
			fmt.Printf("info: 不使用代理访问 %s 成功，响应状态码为 %d\n", website, resp.StatusCode)
			return nil
		}
	}
	return res
}
