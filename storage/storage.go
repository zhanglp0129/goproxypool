package storage

import (
	"github.com/zhanglp0129/goproxypool/config"
	"github.com/zhanglp0129/goproxypool/constant"
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
	InsertProxyAddress(proxyAddress ProxyAddress) error

	// GetAvailableProxyAddress 获取一个可用的代理地址，会提供负载均衡的功能
	GetAvailableProxyAddress(protocol string) (ProxyAddress, error)

	// GetDetectedProxyAddresses 获取待检测的代理地址
	GetDetectedProxyAddresses() ([]ProxyAddress, error)

	// PageProxyAddresses 分页查询代理地址
	PageProxyAddresses(pageNum, pageSize int) (ProxyAddressPageVO, error)

	// UpdateProxyAddress 修改代理地址
	UpdateProxyAddress(proxyAddress ProxyAddress) error

	// DeleteProxyAddress 删除代理地址
	DeleteProxyAddress(id int) error

	// FinishDetection 代理地址检测完成，如果检测失败次数超过阈值，会删除该代理地址
	// accept为是否通过检测
	FinishDetection(id int, accept bool) error

	// FinishUse 完成代理地址的使用
	// success为是否使用成功
	FinishUse(id int, success bool) error
}

// ProxyAddress 代理地址
type ProxyAddress struct {
	ID       int64
	IP       string
	Port     uint16
	Protocol string
}

// ProxyAddressPageVO 代理地址分页查询结果
type ProxyAddressPageVO struct {
	// 合计
	Total int `json:"total"`
	// 待检测数量
	Pends int `json:"pends"`
	// 通过检测数量
	Accepts int `json:"accepts"`
	// 未通过检测数量
	Fails int `json:"Fails"`
	// 结果列表
	Items []struct {
		ID       int                         `json:"id"`
		IP       string                      `json:"ip"`
		Port     uint16                      `json:"port"`
		Protocol string                      `json:"protocol"`
		Status   constant.ProxyAddressStatus `json:"status"`
	} `json:"items"`
}
