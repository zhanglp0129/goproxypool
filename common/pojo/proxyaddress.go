package pojo

import (
	"github.com/zhanglp0129/goproxypool/common/constant"
)

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
	// 正在检测
	Detects int `json:"detects"`
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
