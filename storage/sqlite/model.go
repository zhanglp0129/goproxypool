package sqlite

// StorageModel 持久化存储数据模型
type StorageModel struct {
	ID            int64  `gorm:"primaryKey;autoIncrement"`
	Protocol      string `gorm:"uniqueIndex:idx_address,priority:1;not null;comment:协议"`
	IP            string `gorm:"uniqueIndex:idx_address,priority:2;not null;comment:ip"`
	Port          uint16 `gorm:"uniqueIndex:idx_address,priority:3;not null;comment:端口"`
	AcceptNumber  int    `gorm:"default:0;not null;comment:检测通过次数，负数表示失败次数"`
	EffectiveTime int64  `gorm:"default:0;not null;index;comment:检测生效时间，以完成检测为准，为纳秒级Unix时间戳"`
	UseTime       int64  `gorm:"default:0;not null;index;comment:上次使用时间，以完成使用为准，为纳秒级Unix时间戳"`
}

// TableName 自定义表名
func (StorageModel) TableName() string {
	return "storage"
}
