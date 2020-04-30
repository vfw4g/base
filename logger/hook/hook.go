package hook

import (
	"fmt"
	"github.com/heralight/logrus_mate"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	logrus_mate.RegisterHook("file", NewFileHook)
}

func NewFileHook(options logrus_mate.Options) (hook logrus.Hook, err error) {

	conf := FileLogConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
	}
	w := NewFileWriter(&conf)

	hook = &FileHook{W: w}

	return
}

type FileHook struct {
	W *FileLogWriter
}

func (p *FileHook) Fire(entry *logrus.Entry) (err error) {
	message, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	switch entry.Level {
	case logrus.PanicLevel:
		fallthrough
	case logrus.FatalLevel:
		fallthrough
	case logrus.ErrorLevel:
		return p.W.WriteMsg(message, logrus.ErrorLevel)
	case logrus.WarnLevel:
		return p.W.WriteMsg(message, logrus.WarnLevel)
	case logrus.InfoLevel:
		return p.W.WriteMsg(message, logrus.InfoLevel)
	case logrus.DebugLevel:
		return p.W.WriteMsg(message, logrus.DebugLevel)
	default:
		return nil
	}

	return
}

func (p *FileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////
