package configs

import (
	"github.com/spf13/viper"
	"github.com/vfw4g/base/errors"
)

type Configs struct {
	*viper.Viper
}

// the options will be invorked before read config
type Option func(v *viper.Viper)

// init the config file
// configName sets name for the config file. Does not include extension.
// path	adds the paths to search for the config file in.
func InitConfig(configName string, paths ...string) (cfg *Configs, err error) {
	var vp = viper.New()
	vp.SetConfigName(configName)
	for _, path := range paths {
		vp.AddConfigPath(path)
	}
	if err = vp.ReadInConfig(); err != nil {
		//return nil, errorx.InitializationFailed.WrapWithNoMessage(err)
		return nil, errors.Wrap(err)
	}
	cfg = &Configs{
		vp,
	}
	return
}

func InitConfigLocation(absFile string) (cfg *Configs, err error) {
	var vp = viper.New()
	vp.SetConfigFile(absFile)
	if err = vp.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err)
	}
	cfg = &Configs{
		vp,
	}
	return
}

func InitDefault() (cfg *Configs, err error) {
	defaultDir := "."
	viper.AddConfigPath(defaultDir)
	if err = viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err)
	}
	cfg = &Configs{
		viper.GetViper(),
	}
	return
}
