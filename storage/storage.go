package storage

// Storage 代理地址持久化存储的接口
type Storage interface {
	// InsertProxyAddress 插入一个代理地址。ID应当为0，由系统自动生成
	InsertProxyAddress(proxyAddress ProxyAddress) error

	// GetAvailableProxyAddress 获取一个可用的代理地址，会提供负载均衡的功能。protocol为代理协议
	GetAvailableProxyAddress(protocol string) (ProxyAddress, error)

	// GetDetectedProxyAddresses 获取待检测的代理地址。limit为获取个数上限
	GetDetectedProxyAddresses(limit int) ([]ProxyAddress, error)

	// PageProxyAddresses 分页查询代理地址。分别返回：总记录数、结果列表、可能发生的异常
	PageProxyAddresses(pageNum, pageSize int) (int, []ProxyAddress, error)

	// UpdateProxyAddress 修改代理地址
	UpdateProxyAddress(proxyAddress ProxyAddress) error

	// DeleteProxyAddress 删除代理地址
	DeleteProxyAddress(id int) error

	// FinishDetection 代理地址检测完成，如果检测失败次数超过阈值，会删除该代理地址。accept为是否通过检测
	FinishDetection(id int, accept bool) error
}

// ProxyAddress 代理地址
type ProxyAddress struct {
	ID       int64
	IP       string
	Port     uint16
	Protocol string
}
