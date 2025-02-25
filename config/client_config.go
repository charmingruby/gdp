package config

type ClientConfig struct {
	ServerPort      int `yaml:"server-port"`
	PackageLoadSize int `yaml:"package-load-size"`
}
