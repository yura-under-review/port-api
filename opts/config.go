package opts

import (
	"github.com/spf13/viper"
	"github.com/yura-under-review/port-api/repository"
	"github.com/yura-under-review/port-api/server"
)

const (
	defaultEnvPrefix = "APP"

	defaultLogLevel = "DEBUG"

	keyLogLevel   = "LOG_LEVEL"
	keyPrettyLogs = "PRETTY_LOGS"

	keyServerAddress       = "HTTP_ADDRESS"
	keyServerRootTemplate  = "ROOT_TEMPLATE_FILE"
	keyServerSinkBatchSize = "SINK_BATCH_SIZE"

	keyPortsDomainServiceAddress = "PORTS_DOMAIN_SERVICE_ADDRESS"
)

type Config struct {
	LogLevel   string
	PrettyLogs bool

	PortsDomainService repository.Config
	Server             server.Config
}

func LoadConfigFromEnv() Config {

	viper.AutomaticEnv()

	viper.SetEnvPrefix(defaultEnvPrefix)

	viper.SetDefault(keyLogLevel, defaultLogLevel)
	viper.SetDefault(keyPrettyLogs, false)

	viper.SetDefault(keyServerAddress, ":8080")

	return Config{
		LogLevel:   viper.GetString(keyLogLevel),
		PrettyLogs: viper.GetBool(keyPrettyLogs),

		PortsDomainService: repository.Config{
			Address: viper.GetString(keyPortsDomainServiceAddress),
		},

		Server: server.Config{
			Address:          viper.GetString(keyServerAddress),
			RootTemplateFile: viper.GetString(keyServerRootTemplate),
			SinkBatchSize:    viper.GetInt(keyServerSinkBatchSize),
		},
	}
}
