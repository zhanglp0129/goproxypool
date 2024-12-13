package sqlite

import (
	"errors"
	"github.com/zhanglp0129/goproxypool/common/constant"
	"github.com/zhanglp0129/goproxypool/common/pojo"
	"github.com/zhanglp0129/goproxypool/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math"
	"sync"
	"time"
)

var (
	CFG = config.CFG
)

// InitSqlite 初始化sqlite
func InitSqlite() (*Storage, error) {
	// 获取相关配置
	dsn := CFG.Storage.DSN
	// 连接数据库
	db, err := gorm.Open(sqlite.Open(dsn))
	if err != nil {
		return nil, err
	}

	// 数据库字段迁移
	err = db.AutoMigrate(&StorageModel{})
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

// Storage 持久化存储结构体
type Storage struct {
	db *gorm.DB
	// 每个代理地址的并发数，key: id, value: 并发数
	concurrency sync.Map
	// 代理地址是否正在检测，key：id, value: 是否正在检测
	detecting sync.Map
}

func (s *Storage) InsertProxyAddress(proxyAddress pojo.ProxyAddress) error {
	// 自定义SQL：如果代理地址已存在，则更新其他数据；不存在则插入
	sql := `insert into storage(protocol, ip, port) values (?, ?, ?) 
            on conflict(protocol, ip, port) do
            	update set accept_number=0, effective_time=0, use_time=0`
	return s.db.Raw(sql, proxyAddress.Protocol, proxyAddress.IP, proxyAddress.Port).Error
}

func (s *Storage) GetAvailableProxyAddress(protocol string) (pojo.ProxyAddress, error) {
	maxConcurrency := CFG.Use.MaxConcurrency
	// 获取一个可用的代理地址。未超过检测生效时间，优先选择最久未使用的代理地址。
	model := StorageModel{}
	err := s.db.Where("protocol = ? and effective_time > ? and accept_number > 0", protocol, time.Now().UnixNano()).
		Order("use_time").
		Take(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 无代理地址
		return pojo.ProxyAddress{}, constant.NoProxy
	} else if err != nil {
		// 其他错误
		return pojo.ProxyAddress{}, err
	}

	// 加乐观锁，最多重试3次
	var swapped bool
	for i := 0; i < 3; i++ {
		// 超过最大并发，也认为是无可用的代理地址
		concurrent, _ := s.concurrency.LoadOrStore(model.ID, 1)
		if concurrent.(int) > maxConcurrency {
			return pojo.ProxyAddress{}, constant.NoProxy
		}
		// 获取完成后，并发数+1
		swapped = s.concurrency.CompareAndSwap(model.ID, concurrent, concurrent.(int)+1)
		if swapped {
			break
		}
	}
	if !swapped {
		// 多次尝试乐观锁判断都失效
		return pojo.ProxyAddress{}, constant.NoProxy
	}

	// 构造并返回结果
	return pojo.ProxyAddress{
		ID:       model.ID,
		IP:       model.IP,
		Port:     model.Port,
		Protocol: model.Protocol,
	}, nil
}

func (s *Storage) GetDetectedProxyAddresses() ([]pojo.ProxyAddress, error) {
	// 获取多个待检测代理地址。超过检测生效时间，优先获取最久未检测的代理地址
	limit := CFG.Detect.Number
	var models []StorageModel
	err := s.db.Where("effective_time < ?", time.Now().UnixNano()).
		Order("effective_time").Limit(limit).Find(&models).Error
	if err != nil {
		return nil, err
	}
	// 判断是否正在检测，如果正在检测就删除，用滑动窗口算法优化删除。同时将状态改为正在检测
	var slow int
	for fast := 0; fast < len(models); fast++ {
		model := models[fast]
		s.detecting.LoadOrStore(model.ID, false)
		if s.detecting.CompareAndSwap(model.ID, false, true) {
			// 当前代理地址不在检测中
			models[slow] = models[fast]
			slow++
		}
	}
	models = models[:slow]

	// 构造返回结果
	addresses := make([]pojo.ProxyAddress, 0, len(models))
	for _, model := range models {
		addresses = append(addresses, pojo.ProxyAddress{
			ID:       model.ID,
			IP:       model.IP,
			Port:     model.Port,
			Protocol: model.Protocol,
		})
	}
	return addresses, nil
}

func (s *Storage) PageProxyAddresses(pageNum, pageSize int) (pojo.ProxyAddressPageVO, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) UpdateProxyAddress(proxyAddress pojo.ProxyAddress) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) DeleteProxyAddress(id int) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) FinishDetection(id int64, accept bool) error {
	// 完成代理地址的检测
	// 先获取数据，要求代理地址失效
	var model StorageModel
	err := s.db.Select("id", "accept_number", "effective_time").
		Where("id = ? and effective_time < ?", id, time.Now().UnixNano()).
		Take(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} else if err != nil {
		return err
	}
	// 修改数据库中通过次数和检测生效时间
	if accept {
		// 通过检测
		if model.AcceptNumber > 0 {
			model.AcceptNumber++
		} else {
			model.AcceptNumber = 1
		}
	} else {
		// 未通过检测
		if model.AcceptNumber < 0 {
			model.AcceptNumber++
		} else {
			model.AcceptNumber = -1
		}
	}
	oldEffectiveTime := model.EffectiveTime
	rate := min(CFG.Detect.MaxRate, math.Pow(CFG.Detect.EffectiveRate, math.Abs(float64(model.AcceptNumber))))
	effective := time.Duration(float64(CFG.Detect.EffectiveSeconds)*rate) * time.Second
	model.EffectiveTime = time.Now().UnixNano() + int64(effective)
	// 修改数据，通过代理地址的生效时间加乐观锁
	err = s.db.Select("accept_number", "effective_time").
		Where("id = ? and effective_time = ?", id, oldEffectiveTime).
		Updates(&model).Error
	if err != nil {
		return err
	}

	// 将正在检测改为false
	s.detecting.Store(id, false)
	return nil
}

func (s *Storage) FinishUse(id int64, success bool) error {
	//TODO implement me
	panic("implement me")
}
