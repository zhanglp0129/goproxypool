package utils

import (
	"fmt"
	"github.com/zhanglp0129/goproxypool/common/constant"
	"github.com/zhanglp0129/goproxypool/common/pojo"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// BuildProxyUrl 构建代理url
func BuildProxyUrl(address pojo.ProxyAddress) (*url.URL, error) {
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
