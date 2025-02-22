package config

type ClientConfig struct {
	ServerPort      int                   `yaml:"server-port"`
	PackageLoadSize int                   `yaml:"package-load-size"`
	ClientThreshold ClientThresholdConfig `yaml:"client-thresholds"`
}

type ClientThresholdConfig struct {
	TimeoutInSeconds  int `yaml:"timeout-in-seconds"`
	InitialWindowSize int `yaml:"initial-window-size"`
	MaxWindowSize     int `yaml:"max-window-size"`
	InitialSshthresh  int `yaml:"initial-sshthresh"`
}
