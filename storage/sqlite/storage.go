package sqlite

import (
	"github.com/zhanglp0129/goproxypool/storage"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func (s *Storage) InsertProxyAddress(proxyAddress storage.ProxyAddress) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetAvailableProxyAddress(protocol string) (storage.ProxyAddress, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetDetectedProxyAddresses() ([]storage.ProxyAddress, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) PageProxyAddresses(pageNum, pageSize int) (storage.ProxyAddressPageVO, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) UpdateProxyAddress(proxyAddress storage.ProxyAddress) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) DeleteProxyAddress(id int) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) FinishDetection(id int, accept bool) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) FinishUse(id int, success bool) error {
	//TODO implement me
	panic("implement me")
}
