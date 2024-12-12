package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
)

// 默认配置文件路径
const defaultConfigFilename = "etc/goproxypool.yml"

var CFG *Config

// 初始化配置文件
func init() {
	// 获取配置文件路径
	flag.Parse()
	configFilename := flag.String("c", defaultConfigFilename, "配置文件路径")
	// 读取配置文件内容
	content, err := os.ReadFile(*configFilename)
	if err != nil {
		panic(err)
	}
	// 解析配置文件
	CFG = new(Config)
	err = yaml.Unmarshal(content, CFG)
}

// Config 配置结构体
type Config struct {
	Proxy   ProxyConfig
	Panel   PanelConfig
	Detect  DetectConfig
	Use     UseConfig
	Storage StorageConfig
	Log     LogConfig
}

// ProxyConfig 代理相关配置
type ProxyConfig struct {
	Http HttpProxyConfig
}

// HttpProxyConfig http代理相关配置
type HttpProxyConfig struct {
	IP      string
	Port    uint16
	NoProxy string `yaml:"no_proxy"`
}

// PanelConfig 管理面板相关配置
type PanelConfig struct {
	IP   string
	Port uint16
}

// DetectConfig 可用性检测相关配置
type DetectConfig struct {
	Number           int
	Interval         int
	Timeout          int
	EffectiveSeconds int     `yaml:"effective_seconds"`
	EffectiveRate    float64 `yaml:"effective_rate"`
	MaxRate          float64 `yaml:"max_rate"`
	Attempts         int
	Websites         []string
	DirectInterval   int `yaml:"direct_interval"`
}

// UseConfig 使用代理地址相关配置
type UseConfig struct {
	Detect    bool
	Timeout   int
	Attempts  int
	Addresses int
}

// StorageConfig 持久化存储相关配置
type StorageConfig struct {
	Type string
	DSN  string
}

// LogConfig 日志相关配置
type LogConfig struct {
	Level string
	File  string
}
