package logger

import (
	"github.com/heirko/go-contrib/logrusHelper"
	"github.com/sirupsen/logrus"
	"github.com/vfw4g/base/configs"
	"github.com/vfw4g/base/errors"
	_ "github.com/vfw4g/base/logger/hook"
)

///////////////////////////////////////////////////////////////////////////////////////////////////
type Logger struct {
	*logrus.Logger
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// the options will be invorked before read config
type Option func(l *logrus.Logger)

var REPORT_CALLER_OPTION Option = func(l *logrus.Logger) {
	l.SetReportCaller(true)
}

///////////////////////////////////////////////////////////////////////////////////////////////////
func InitLogger(lgLocation string, ops ...Option) (vfwlogger *Logger, err error) {
	var cfg *configs.Configs
	if cfg, err = configs.InitConfigLocation(lgLocation); err != nil {
		return nil, errors.Wrap(err)
	}
	lg := logrus.New()

	for _, op := range ops {
		op(lg)
	}

	// Read configuration
	//mate.RegisterWriter("rotatelogs", NewRotatelogsWriter)
	var c = logrusHelper.UnmarshalConfiguration(cfg.Viper) // Unmarshal configuration from Viper
	if err = logrusHelper.SetConfig(lg, c); err != nil {   // for e.g. apply it to logrus default instance
		return nil, errors.Wrap(err)
	}
	vfwlogger = &Logger{lg}
	return
	// ### End Read Configuration
}

///////////////////////////////////////////////////////////////////////////////////////////////////

func (l *Logger) Errorvn(arg interface{}) {
	l.Errorf("%+v\n", arg)
}

func (l *Logger) Fatalvn(arg interface{}) {
	l.Fatalf("%+v\n", arg)
}
