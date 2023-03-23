package conf

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

const (
	AppName = "messenger"

	ParamPort     = "bind.port"
	ParamBindHost = "bind.host"

	ParamOpenAPIKey = "openai.key"

	ParamDebug = "debug"

	ParamMongoDB = "storage.mongodb.uri"

	ParamMessageRetention = "chat.retention"

	ParamTrustedProxies = "server.trusted-proxies"
)

func ParseConfig() {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/" + AppName + "/") // path to look for the config file in
	viper.AddConfigPath("$HOME/." + AppName)     // call multiple times to add many search paths
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalln("Cannot find config.yaml")
		} else {
			log.Fatalf("Cannot read config file: %v", err)
		}
	}
}

func BindPort() int {
	return viper.GetInt(ParamPort)
}

func BindHost() string {
	return viper.GetString(ParamBindHost)
}

func OpenAPIKey() string {
	return viper.GetString(ParamOpenAPIKey)
}

func Debug() bool {
	return viper.GetBool(ParamDebug)
}

func MongoURI() string {
	return viper.GetString(ParamMongoDB)
}

func TrustedProxies() []string {
	return viper.GetStringSlice(ParamTrustedProxies)
}

func Retention() time.Duration {
	return time.Duration(viper.GetInt64(ParamMessageRetention)) * time.Second
}
