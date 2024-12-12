package sqlite

import "time"

// StorageModel 持久化存储数据模型
type StorageModel struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	Protocol     string    `gorm:"uniqueIndex:idx_address"`
	IP           string    `gorm:"uniqueIndex:idx_address"`
	Port         uint16    `gorm:"uniqueIndex:idx_address"`
	AcceptNumber int       `gorm:"default:0;comment:检测通过次数，负数表示失败次数"`
	DetectTime   time.Time `gorm:"index;comment:上次检测时间"`
	UseTime      time.Time `gorm:"index;comment:上次使用时间"`
}

// TableName 自定义表名
func (StorageModel) TableName() string {
	return "storage"
}
