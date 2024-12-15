package detector

import (
	"errors"
	"fmt"
	"github.com/zhanglp0129/goproxypool/common/constant"
	"github.com/zhanglp0129/goproxypool/common/pojo"
	"github.com/zhanglp0129/goproxypool/utils"
	"math/rand"
	"net/http"
	"time"
)

// 运行代理地址的可用性检测
func runAddressDetect() {
	// 启动时就执行一次检测
	doDetect()
	// 开启一个计时器，每隔一段时间检测一次
	interval := time.Duration(CFG.Detect.Interval) * time.Second
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		doDetect()
	}
}

// 执行可用性检测
func doDetect() {
	// 获取指定个数的代理地址
	addresses, err := Storage.GetDetectedProxyAddresses()
	if err != nil {
		// TODO 记录日志
		fmt.Printf("error: 获取待检测代理地址错误 %v\n", err)
		return
	}
	// 开启goroutine，执行检测
	for _, address := range addresses {
		go Detect(address, true)
	}
}

// Detect 检测代理地址的可用性。acceptFinish通过时是否完成检测
func Detect(address pojo.ProxyAddress, acceptFinish bool) error {
	// 获取重试次数
	attempts := CFG.Detect.Attempts
	// 返回的error
	var res error
	// 判断是否为无检测网站
	noWebsite := true
	for i := 0; i < attempts; i++ {
		// 获取检测代理地址使用的网站
		var website string
		err := func() error {
			websitesMutex.RLock()
			defer websitesMutex.RUnlock()
			if len(availableWebsites) == 0 {
				// TODO 记录日志
				fmt.Printf("warn: 无可用的检测网站\n")
				return constant.NoDetectWebsite
			}
			website = availableWebsites[rand.Intn(len(availableWebsites))]
			return nil
		}()
		timeout := time.Duration(CFG.Detect.Timeout) * time.Second
		if err != nil {
			// 获取检测网站失败，按照超时处理
			time.Sleep(timeout)
			continue
		} else {
			// TODO 打印日志
			fmt.Printf("info: 检测代理地址 %v 使用的网站为 %s\n", address, website)
		}
		noWebsite = false

		// 构建代理对象
		proxyUrl, err := utils.BuildProxyUrl(address)
		if err != nil {
			// TODO 记录日志
			fmt.Printf("error: 构建代理对象错误，代理地址为 %v\n", address)
			return err
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		client := http.Client{
			Transport: transport,
			Timeout:   timeout,
		}

		// 向网站发送代理请求
		resp, err := client.Get(website)
		if err != nil {
			// TODO 记录 info 日志
			fmt.Printf("info: 检测代理地址 %v 错误 %v\n", address, err)
			res = errors.Join(res, err)
		} else if resp.StatusCode != 200 {
			// TODO 记录 info 日志
			fmt.Printf("info: 检测代理地址 %v 错误 %s\n", address, resp.Status)
			res = errors.Join(res, errors.Join(errors.New(resp.Status)))
		} else {
			// TODO 记录 info 日志
			fmt.Printf("info: 使用代理 %v 访问 %s 成功，响应状态码为 %d\n", address, website, resp.StatusCode)
			// 检测完成
			if !acceptFinish {
				return nil
			}
			err = Storage.FinishDetection(address.ID, true)
			if err != nil {
				// TODO 记录日志
				fmt.Printf("error: 完成代理地址 %v 检测错误 %v\n", address, err)
				return constant.FinishDetectError
			}
			return nil
		}
	}
	if noWebsite {
		return constant.NoDetectWebsite
	}
	// 检测完成
	err := Storage.FinishDetection(address.ID, false)
	if err != nil {
		// TODO 记录日志
		fmt.Printf("error: 完成代理地址 %v 检测错误 %v\n", address, err)
		res = errors.Join(res, constant.FinishDetectError)
	}
	return res
}
