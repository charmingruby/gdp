package config

type ServerConfig struct {
	Port                int                       `yaml:"port"`
	CongestionThreshold CongestionThresholdConfig `yaml:"congestion-thresholds"`
}

type CongestionThresholdConfig struct {
	PackageLoss float32 `yaml:"package-loss"`
}
