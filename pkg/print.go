package pkg

import (
	"github.com/sirupsen/logrus"
)

func InfoPrint(layer, status, message string) {
	logrus.WithFields(logrus.Fields{"layer": layer, "status": status}).Info(message)
}
func ErrPrint(layer, status, message string, err error) {
	logrus.WithFields(logrus.Fields{"layer": layer, "status": status, "error": err}).Error(message)
}
