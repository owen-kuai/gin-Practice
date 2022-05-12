package conf

import (
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

var (
	config     *Config
	configInit sync.Once
)

func InitConfig(configDir, appName, fileType string) {
	configInit.Do(func() {
		if config == nil {
			config = &Config{}
		}
		if err := InitVipper(configDir, appName, fileType); err != nil {
			panic(err)
		}
		return
	})

}

//InitVipper init config setting for viper , set config file name , config dir and config type
// supported file type "json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "tfvars", "dotenv", "env", "ini"
func InitVipper(configDir, appName, fileType string) error {

	viper.SetConfigName(appName)
	viper.SetConfigType(fileType)
	viper.AddConfigPath(configDir)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	if err := viper.Unmarshal(config); err != nil {
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
	return nil
}

//type Config interface {
//	//ImplementationWithDefault complete the config with default value
//	ImplementationWithDefault() Config
//}
type Config struct {
	LogConf        LogConf `mapstructure:"logger"`
	Host           string  `mapstructure:"host"`
	Port           int     `mapstructure:"port"`
	EnableSwagger  bool    `mapstructure:"enable_swagger"`
	KubeConfig     string  `mapstructure:"kube_config"`
	ClientTimeout  int     `mapstructure:"client_timeout"`
	OckleUrl       string  `mapstructure:"ockle_url"`
	Gateway        string  `mapstructure:"gateway"`
	WalmUrl        string  `mapstructure:"walm_url"`
	TmpDir         string  `mapstructure:"tmp_dir"`
	HamurapiUrl    string  `mapstructure:"hamurapi_url"`
	ProductMetaDir string  `mapstructure:"product_meta_dir"`
	Version        string
}

type LogConf struct {
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxAge     int    `mapstructure:"maxAge"`
	MaxBackups int    `mapstructure:"maxBackups"`
	Compress   bool   `mapstructure:"compress"`
}

func GetConfig() Config {
	return *config
}
