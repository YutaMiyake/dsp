package common

import (
	l "github.com/sirupsen/logrus"
	"os"
	"sync"
)

const (
	logfile = "output.log"
)

var (
	loggerOnce = new(sync.Once)
	Logger     *l.Logger // singleton
)

func SetupLogger() {
	loggerOnce.Do(func() {
		log := l.New()
		log.Formatter = new(l.TextFormatter)
		log.Level = l.InfoLevel
		log.Out = os.Stdout
		Logger = log
	})
}
