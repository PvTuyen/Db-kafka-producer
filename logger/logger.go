package logger

import (
	log "github.com/sirupsen/logrus"
)

func InitLog()  {
	log.SetFormatter(&log.JSONFormatter{})
}
func LogInfo(msg string, detail string)  {
	log.WithFields(
		log.Fields{
			"detail": detail,
		},
	).Info(msg)
}
func LogError(msg string, detail string)  {
	log.WithFields(
		log.Fields{
			"detail": detail,
		},
	).Error(msg)
}