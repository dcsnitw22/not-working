package logger

import "github.com/sirupsen/logrus"

var log *logrus.Logger = logrus.New()
var AppLog *logrus.Entry = log.WithFields(logrus.Fields{"AMF": "app"})
