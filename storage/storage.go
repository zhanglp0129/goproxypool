package storage

// Storage 代理地址持久化存储的接口
type Storage interface {
	// InsertProxyAddress 插入一个代理地址。当ID为0时，应当自动生成id
	InsertProxyAddress(proxyAddress ProxyAddress) error

	// RandomAvailableProxyAddresses 随机获取可用的代理地址，limit为获取个数上限
	RandomAvailableProxyAddresses(limit int) ([]ProxyAddress, error)

	// RandomDetectedProxyAddresses 随机获取待检测的代理地址，limit为获取个数上限
	RandomDetectedProxyAddresses(limit int) ([]ProxyAddress, error)

	// PageProxyAddresses 分页查询代理地址。分别返回：总记录数、结果列表、可能发生的异常
	PageProxyAddresses(pageNum, pageSize int) (int, []ProxyAddress, error)

	// UpdateProxyAddress 修改代理地址
	UpdateProxyAddress(proxyAddress ProxyAddress) error

	// DeleteProxyAddress 删除代理地址
	DeleteProxyAddress(id int) error

	// FinishDetection 代理地址检测完成。success为是否检测成功
	FinishDetection(id int, success bool) error
}

// ProxyAddress 代理地址
type ProxyAddress struct {
	ID       int64
	IP       string
	Port     uint16
	Protocol string
}
