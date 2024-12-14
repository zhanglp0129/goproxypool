package detector

import (
	"errors"
	"fmt"
	"github.com/zhanglp0129/goproxypool/common/constant"
	"github.com/zhanglp0129/goproxypool/common/pojo"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
			fmt.Printf("error: 获取待检测代理地址错误 %v\n", err)
			continue
		}
		// 开启goroutine，执行检测
		for _, address := range addresses {
			go func() {
				err = detect(address)
				if err != nil {
					// TODO 记录 info 日志
					fmt.Printf("info: 检测代理地址 %v 错误 %v\n", address, err)
					return
				}
				// 检测完成
				err = Storage.FinishDetection(address.ID, err == nil)
				if err != nil {
					// TODO 记录日志
					fmt.Printf("error: 完成代理地址 %v 检测错误 %v\n", address, err)
					return
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
		// 获取检测代理地址使用的网站
		var website string
		err := func() error {
			websitesMutex.RLock()
			defer websitesMutex.RUnlock()
			if len(availableWebsites) == 0 {
				// TODO 记录日志
				fmt.Printf("error: 无可用的检测网站\n")
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
		}

		// 构建代理对象
		proxyUrl, err := buildProxyUrl(address)
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
			res = errors.Join(res, err)
		} else {
			// TODO 记录 info 日志
			fmt.Printf("info: 使用代理 %v 访问 %s 成功，响应状态码为 %d\n", address, website, resp.StatusCode)
			return nil
		}
	}
	return res
}

// 构建代理url
func buildProxyUrl(address pojo.ProxyAddress) (*url.URL, error) {
	builder := strings.Builder{}
	// 构建协议
	builder.WriteString(address.Protocol)
	builder.WriteString("://")
	// 构建ip
	ip := net.ParseIP(address.IP)
	if ip == nil {
		return nil, constant.IPFormatError
	}
	if ip.To4() != nil {
		// ipv4
		builder.WriteString(ip.String())
	} else {
		// ipv6
		builder.WriteRune('[')
		builder.WriteString(ip.String())
		builder.WriteRune(']')
	}
	// 构建端口
	builder.WriteRune(':')
	builder.WriteString(strconv.Itoa(int(address.Port)))
	rawUrl := builder.String()
	// TODO 打印 debug 日志
	fmt.Printf("debug: 代理url为 %s\n", rawUrl)
	return url.Parse(rawUrl)
}
