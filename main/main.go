package main

import (
	"github.com/zhanglp0129/goproxypool/detector"
	"github.com/zhanglp0129/goproxypool/proxy"
)

func main() {
	detector.Run()
	proxy.Run()
	select {}
}
