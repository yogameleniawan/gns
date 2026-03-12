package config

type App struct {
	Name             string `yaml:"name"`
	Version          string `yaml:"version"`
	Description      string `yaml:"description"`
	Port             string `yaml:"port"`
	RequestPerSecond int    `yaml:"request_per_second"`
	Burst            int    `yaml:"burst"`
}
