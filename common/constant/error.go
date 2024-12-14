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
)
