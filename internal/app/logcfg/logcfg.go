package logcfg

import (
	"fmt"
	"io"

	// "log/slog"
	"os"
	"path"
	"runtime"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func RunLoggerConfig(EnvLogs string) {

	// TODO
	// var hlg *slog.Logger
	// slog.SetDefault(hlg)
	// slog.SetLogLoggerLevel(slog.LevelInfo)

	logLevel, err := logrus.ParseLevel(EnvLogs)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)

	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			_, filename := path.Split(f.File)
			filename = fmt.Sprintf("%s.%d.%s", filename, f.Line, f.Function)
			return "", filename
		},
	})
	mw := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "gophermart.log",
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     30,
	})
	logrus.SetOutput(mw)
	// slog.NewTextHandler(mw, &slog.HandlerOptions{})
}
