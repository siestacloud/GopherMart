package pkg

import (
	"github.com/sirupsen/logrus"
)

func InfoPrint(layer, status string, message ...interface{}) {
	logrus.WithFields(logrus.Fields{"layer": layer, "status": status}).Info(message...)
}

func WarnPrint(layer, status string, message ...interface{}) {
	logrus.WithFields(logrus.Fields{"layer": layer, "status": status}).Warn(message...)
}

func ErrPrint(layer, status interface{}, message ...interface{}) {
	logrus.WithFields(logrus.Fields{"layer": layer, "status": status}).Error(message...)
}
