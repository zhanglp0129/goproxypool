package config

type Config struct {
	Proxy   ProxyConfig
	Panel   PanelConfig
	Detect  DetectConfig
	Storage StorageConfig
	Log     LogConfig
}

type ProxyConfig struct {
	Http HttpProxyConfig
}

type HttpProxyConfig struct {
	IP   string
	Port uint16
	NoIP string
}

type PanelConfig struct {
	IP   string
	Port uint16
}

type DetectConfig struct {
	Number           int
	Interval         int
	EffectiveSeconds int
	DeleteThreshold  int
	Websites         []string
	DirectInterval   int
}

type StorageConfig struct {
	Type string
	DSN  string
}

type LogConfig struct {
	Level string
	File  string
}
