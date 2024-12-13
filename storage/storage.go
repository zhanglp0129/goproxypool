package storage

import (
	"github.com/zhanglp0129/goproxypool/common/pojo"
	"github.com/zhanglp0129/goproxypool/config"
	"github.com/zhanglp0129/goproxypool/storage/sqlite"
)

var (
	Storage IStorage
	CFG     = config.CFG
)

const (
	Sqlite = "sqlite"
)

// 初始化持久化存储实例
func init() {
	// 获取持久化存储相关配置
	typ := CFG.Storage.Type
	switch typ {
	case Sqlite:
		s, err := sqlite.InitSqlite()
		if err != nil {
			panic(err)
		}
		Storage = s
	default:
		panic("持久化类型不合法")
	}
}

// IStorage 代理地址持久化存储的接口
type IStorage interface {
	// InsertProxyAddress 插入一个代理地址
	// proxyAddress.ID无需指定，由系统自动生成。如果该代理地址已存在，则刷新
	InsertProxyAddress(proxyAddress pojo.ProxyAddress) error

	// GetAvailableProxyAddress 获取一个可用的代理地址，会提供负载均衡的功能
	GetAvailableProxyAddress(protocol string) (pojo.ProxyAddress, error)

	// GetDetectedProxyAddresses 获取待检测的代理地址
	GetDetectedProxyAddresses() ([]pojo.ProxyAddress, error)

	// PageProxyAddresses 分页查询代理地址
	PageProxyAddresses(pageNum, pageSize int) (pojo.ProxyAddressPageVO, error)

	// UpdateProxyAddress 修改代理地址
	UpdateProxyAddress(proxyAddress pojo.ProxyAddress) error

	// DeleteProxyAddress 删除代理地址
	DeleteProxyAddress(id int) error

	// FinishDetection 代理地址检测完成
	// accept为是否通过检测
	FinishDetection(id int64, accept bool) error

	// FinishUse 完成代理地址的使用
	// success为是否使用成功
	FinishUse(id int64, success bool) error
}
