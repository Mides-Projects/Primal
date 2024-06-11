package config

type Yaml struct {
	MongoUri string `yaml:"mongodb-uri"`
	RedisUri string `yaml:"redis-uri"`
	Port     string `yaml:"port"`
}
