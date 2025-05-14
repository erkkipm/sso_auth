package configs

import "time"

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	NameProject string `yaml:"name_project" env:"NAME_PROJECT" env-required:"true"`
	//MaxConcurrentConnections int64    `yaml:"max_concurrent_connections" env:"MAX_CONCURRENT_CONNECTIONS" env-default:"10"`
	//Email                    string   `yaml:"email" env:"EMAIL" env-default:"erkkipm@er-company.eu"`
	//Host                     string   `yaml:"host" env:"HTTP-HOST" env-default:"localhost"`
	//Domains                  []string `yaml:"domains" env:"DOMAINS" env-required:"true"`
	GRPC struct {
		//Host         string        `yaml:"ip" env:"GRPC-IP" env-default:"217.0.0.1"`
		Port string `yaml:"port" env:"GRPC-PORT" env-default:"8080"`
		//PortSSL      int           `yaml:"port-ssl" env:"GRPC-PORT-SSL" env-default:"8443"`
		ReadTimeout  time.Duration `yaml:"read-timeout" env:"GRPC-READ-TIMEOUT" env-default:"5s"`
		WriteTimeout time.Duration `yaml:"write-timeout" env:"GRPC-WRITE-TIMEOUT" env-default:"5s"`
	} `yaml:"grpc"`
	TokenTTL int64   `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"3600"`
	MongoDB  MongoDB `yaml:"mongodb"`
	JWTKey   string  `yaml:"jwt_key" env:"JWT_KEY" env-required:"true" env-default:"ERKKIPM"`
}

type MongoDB struct {
	Host       string            `yaml:"host" env:"MONGODB_HOST" env-required:"true" env-default:"localhost"`
	Port       string            `yaml:"port" env:"MONGODB_PORT" env-required:"true" env-default:"27017"`
	Username   string            `yaml:"username" env:"MONGODB_USERNAME" env-required:"true"`
	Password   string            `yaml:"password" env:"MONGODB_PASSWORD" env-required:"true"`
	Database   string            `yaml:"database" env:"MONGODB_DATABASE" env-required:"true" env-default:"main"`
	AuthDB     string            `yaml:"auth_db" env:"MONGODB_AUTH_DB" env-required:"true" env-default:"admin"`
	Collection MongoDBCollection `yaml:"mongodb_collection"`
}

type MongoDBCollection struct {
	Users string `yaml:"users" env:"MONGODB_COLLECTION_USERS" env-required:"true" env-default:"users"`
}
