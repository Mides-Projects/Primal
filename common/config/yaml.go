package config

type Yaml struct {
	MongoUri     string `yaml:"mongodb-uri"`
	MongoDBName  string `yaml:"mongodb-name"`
	RedisUri     string `yaml:"redis-uri"`
	RedisChannel string `yaml:"redis-channel"`
	Port         string `yaml:"port"`
	Key          string `yaml:"key"`
}
