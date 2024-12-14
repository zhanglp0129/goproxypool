package main

import (
	"github.com/zhanglp0129/goproxypool/detector"
	_ "github.com/zhanglp0129/goproxypool/storage"
)

func main() {
	detector.Run()
	select {}
}
