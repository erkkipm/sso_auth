package configs

type Config struct {
	Port     string `yaml:"port"`
	MongoURI string `yaml:"mongo_uri"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
}
