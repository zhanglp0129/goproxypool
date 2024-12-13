package constant

// GoproxypoolError 自定义异常
type GoproxypoolError string

func (e GoproxypoolError) Error() string {
	return string(e)
}

const (
	// NoProxy 无代理地址异常
	NoProxy GoproxypoolError = "no proxy address error"
)
