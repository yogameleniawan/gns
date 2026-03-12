package config

type Cache struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}
