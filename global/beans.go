package global

import (
	"github.com/vfw4g/base/bean"
	"github.com/vfw4g/base/configs"
	"github.com/vfw4g/base/logger"
)

var (
	Logger logger.Logger
	Config configs.Configs
)

func init() {
	bean.GetBeanByNameDelay("github.com/vfw4g/base/logger.Logger", &Logger)
	bean.GetBeanByNameDelay("github.com/vfw4g/base/configs.Configs", &Config)
}
