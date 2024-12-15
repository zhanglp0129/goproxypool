package constant

// GoproxypoolError 自定义异常
type GoproxypoolError string

func (e GoproxypoolError) Error() string {
	return string(e)
}

const (
	// NoProxy 无代理地址异常
	NoProxy GoproxypoolError = "no proxy address error"
	// UseProxyError 使用代理地址失败
	UseProxyError GoproxypoolError = "use proxy address error"
	// IPFormatError IP地址格式错误
	IPFormatError GoproxypoolError = "ip address format error"
	// NoDetectWebsite 无检测代理地址使用的网站
	NoDetectWebsite GoproxypoolError = "no detect website error"
	// FinishDetectError 完成代理地址检测错误
	FinishDetectError GoproxypoolError = "finish detect address error"
)
