package sqlite

import (
	"github.com/zhanglp0129/goproxypool/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var CFG = config.CFG

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
		db,
	}, nil
}
