package config

import (
	"fmt"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type constantViper struct {
	State       string
	MongoConfig struct {
		Address                  []string `mapstructure:"address"`
		Timeout                  string   `mapstructure:"timeout"`
		Username                 string   `mapstructure:"username"`
		Password                 string   `mapstructure:"password"`
		MaxConnectionIdleTimeOut string   `mapstructure:"max_connection_idle_timeout"`
		AuthSource               string   `mapstructure:"auth_source"`
		ReplicaSetName           string   `mapstructure:"replica_set_name"`
		MaxPoolSize              int      `mapstructure:"max_pool_size"`
		DatabaseLocation         string   `mapstructure:"database_location"`
		QueryMaxLimit            int      `mapstructure:"query_max_limit"`
	} `mapstructure:"mongo_config"`
	Collection struct {
		UserLocation string `mapstructure:"user_location"`
	} `mapstructure:"collection_name"`
}

var CV = &constantViper{}

const (
	stageLocal = "local"
	stageUAT   = "uat"
	stagePrd   = "prd"
)

func InitConfig(configPath string, stage string) error {
	curState := "local"
	switch stage {
	case "local", "localhost", "l":
		curState = stageLocal
	case "prod", "production", "p", "prd":
		curState = stagePrd
	case "uat":
		curState = stageUAT
	default:
		curState = stageLocal
	}
	log.Infoln("running on state: ", curState)
	configFileName := fmt.Sprintf("config-%s.yml", curState)
	configFilePath := path.Join(configPath, configFileName)
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("invalid read in config: %w", err)
	}
	loadConfig := constantViper{}
	if err := viper.Unmarshal(&loadConfig); err != nil {
		return fmt.Errorf("invalid unmarshal config: %s", err)
	}
	loadConfig.State = curState
	CV = &loadConfig
	return nil
}
