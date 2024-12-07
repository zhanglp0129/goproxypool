package sqlite

import (
	"github.com/zhanglp0129/goproxypool/storage"
)

// NewSqlite 创建一个sqlite持久化存储实例
func NewSqlite(dsn string) (*Storage, error) {
	return &Storage{}, nil
}

// Storage sqlite持久化存储
type Storage struct {
}

func (s Storage) InsertProxyAddress(proxyAddress storage.ProxyAddress) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetAvailableProxyAddress(protocol string) (storage.ProxyAddress, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetDetectedProxyAddresses(limit int) ([]storage.ProxyAddress, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) PageProxyAddresses(pageNum, pageSize int) (int, []storage.ProxyAddress, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) UpdateProxyAddress(proxyAddress storage.ProxyAddress) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) DeleteProxyAddress(id int) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) FinishDetection(id int, accept bool) error {
	//TODO implement me
	panic("implement me")
}
