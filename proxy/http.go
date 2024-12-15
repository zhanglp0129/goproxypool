package proxy

import (
	"errors"
	"fmt"
	"github.com/elazarl/goproxy"
	"github.com/zhanglp0129/goproxypool/common/constant"
	"github.com/zhanglp0129/goproxypool/common/pojo"
	"github.com/zhanglp0129/goproxypool/detector"
	"github.com/zhanglp0129/goproxypool/utils"
	"net/http"
	"time"
)

// 运行http代理服务器
func runHttp() {
	proxy := goproxy.NewProxyHttpServer()
	// TODO 修改调试日志
	proxy.Verbose = true
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(interceptHttpProxy)
	address := fmt.Sprintf("%s:%d", CFG.Proxy.Http.IP, CFG.Proxy.Http.Port)
	// TODO 打印 debug 日志
	fmt.Printf("debug: http代理服务器监听 %s\n", address)
	err := http.ListenAndServe(address, proxy)
	if err != nil {
		// TODO 记录日志
		fmt.Printf("error: 运行http代理 %s 错误 %v\n", address, err)
		panic(err)
	}
}

// 拦截http代理
func interceptHttpProxy(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// TODO 打印 debug 日志
	fmt.Printf("debug: 拦截代理请求 %v\n", req)
	// 获取每次使用代理地址的数量
	addresses := CFG.Use.Addresses
	// 备份requestURI
	requestURI := req.RequestURI
	req.RequestURI = ""
	for i := 0; i < addresses; i++ {
		// 获取一个可用的代理地址
		proxyAddress, err := Storage.GetAvailableProxyAddress(constant.Http)
		if err != nil {
			// TODO 记录日志，无可用的代理地址
			fmt.Printf("warn: 无可用的代理地址 %v\n", err)
			waiting := time.Duration(CFG.Use.NoProxyWaiting) * time.Second
			time.Sleep(waiting)
			continue
		}
		// 使用http代理地址
		resp, err := useHttpProxyAddress(req, proxyAddress)
		if err != nil {
			// 使用代理地址失败，执行一次检测
			failDetect := CFG.Use.FailDetect
			if failDetect {
				go detector.Detect(proxyAddress)
			}
		} else {
			return req, resp
		}
	}

	// 无可用的代理地址，获取相应配置，并根据配置执行相应策略
	noProxyPolicy := CFG.Use.NoProxy
	switch noProxyPolicy {
	case constant.ErrorPolicy:
		// 错误策略
		return req, goproxy.NewResponse(req, goproxy.ContentTypeText,
			http.StatusInternalServerError, "no available proxy address",
		)
	case constant.DirectPolicy:
		// 直连策略
		req.RequestURI = requestURI
		return req, nil
	default:
		// TODO 记录日志，无可用代理执行策略错误
		fmt.Printf("error: 不支持该策略 %s\n", noProxyPolicy)
		panic("不支持该策略 %s")
	}
}

// 使用http代理地址
func useHttpProxyAddress(req *http.Request, proxyAddress pojo.ProxyAddress) (*http.Response, error) {
	attempts := CFG.Use.Attempts
	var res error
	for i := 0; i < attempts; i++ {
		// 发送代理请求
		req.RequestURI = ""
		resp, err := reqProxy(proxyAddress, req)
		if err != nil {
			// TODO 记录日志
			fmt.Printf("warn: 使用代理地址 %v 失败 %d 次 %v\n", proxyAddress, i+1, err)
			res = errors.Join(res, err)
		} else {
			// TODO 记录日志
			fmt.Printf("info: 使用代理 %v 成功，响应状态码为 %d\n", proxyAddress, resp.StatusCode)
			err = Storage.FinishUse(proxyAddress.ID, true)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}
	}
	err := Storage.FinishUse(proxyAddress.ID, false)
	if err != nil {
		return nil, err
	}
	return nil, res
}

// 发送代理请求
func reqProxy(proxyAddress pojo.ProxyAddress, req *http.Request) (*http.Response, error) {
	// 构建代理url
	url, err := utils.BuildProxyUrl(proxyAddress)
	if err != nil {
		// TODO 记录日志
		fmt.Printf("error: 构建代理url错误 %v\n", err)
		return nil, err
	}

	// 构建代理请求
	transport := &http.Transport{
		Proxy: http.ProxyURL(url),
	}
	timeout := time.Duration(CFG.Use.Timeout) * time.Second
	client := http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	// 发送请求
	return client.Do(req)
}
