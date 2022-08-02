package configuration

import (
	"github.com/spf13/viper"
	"strings"
)

func NewConfig() *Configurations {
	var (
		config   *Configurations
		name     = "configuration"
		typeName = "yaml"
		path     = "./src/configuration"
		err      error
	)

	viper.SetConfigName(name)
	viper.SetConfigType(typeName)
	viper.AddConfigPath(path)

	if err = viper.ReadInConfig(); err != nil {
		return nil
	}

	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err = viper.Unmarshal(&config); err != nil {
		return nil
	}

	return config
}

type Configurations struct {
	Database DatabaseConfigurations
}

type DatabaseConfigurations struct {
	ConnectionString string `mapstructure:"connection_string"`
}
