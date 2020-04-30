package hook

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"time"
)

///////////////////////////////////////////////////////////////////////////////////////////////////
type FileLogConfig struct {
	GlobPattern   string       `json:"glob-pattern"`
	LinkName      string       `json:"link-name"`
	MaxAge        string       `json:"max-age"`
	RotationTime  string       `json:"rotation-time"`
	Clock         string       `json:"clock"`
	RotationCount int          `json:"rotation-count"`
	Level         logrus.Level `json:"level"`
}

func parseDurationPrintErr(s string) (d time.Duration) {
	if d, err := time.ParseDuration(s); err != nil {
		fmt.Println(err)
		return d
	} else {
		return d
	}
}

func (c FileLogConfig) GetMaxAge() (d time.Duration) {
	return parseDurationPrintErr(c.MaxAge)
}

func (c FileLogConfig) GetRotationTime() (d time.Duration) {
	return parseDurationPrintErr(c.RotationTime)
}

func (c FileLogConfig) GetClock() rotatelogs.Clock {
	if c.Clock == "UTC" {
		return rotatelogs.UTC
	} else {
		return rotatelogs.Local
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// FileLogWriter implements LoggerInterface.
// It writes messages by lines limit, file size limit, or time frequency.
type FileLogWriter struct {
	//*log.Logger
	Rl *rotatelogs.RotateLogs
	C  *FileLogConfig
}

// create a FileLogWriter returning as LoggerInterface.
func NewFileWriter(conf *FileLogConfig) *FileLogWriter {
	w := &FileLogWriter{
		C: conf,
	}
	rl, _ := rotatelogs.New(conf.GlobPattern,
		rotatelogs.WithLinkName(conf.LinkName),              // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(conf.GetMaxAge()),             // 文件最大保存时间
		rotatelogs.WithRotationTime(conf.GetRotationTime()), // 日志切割时间间隔
		rotatelogs.WithClock(conf.GetClock()),
		rotatelogs.WithRotationCount(uint(conf.RotationCount)),
	)
	w.Rl = rl
	return w
}

// write logger message into file.
func (w *FileLogWriter) WriteMsg(msg string, level logrus.Level) error {
	if level > w.C.Level {
		return nil
	}
	w.Rl.Write([]byte(msg))
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
