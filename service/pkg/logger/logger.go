package logger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func Init() *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	return &Logger{log}
}

func (l *Logger) withContext(c *fiber.Ctx) *logrus.Entry {
	return l.WithFields(logrus.Fields{
		"method": c.Method(),
		"path":   c.Path(),
	})
}

func (l *Logger) InfoCtx(c *fiber.Ctx, message string) {
	l.withContext(c).Info(message)
}

func (l *Logger) WarnCtx(c *fiber.Ctx, message string) {
	l.withContext(c).Warn(message)
}

func (l *Logger) ErrorCtx(c *fiber.Ctx, message string, err error) {
	l.withContext(c).WithField("error", err.Error()).Error(message)
}

func (l *Logger) WarnWithError(c *fiber.Ctx, message string, err error) {
	l.withContext(c).WithField("error", err).Warn(message)
}

func (l *Logger) ErrorWithFields(c *fiber.Ctx, message string, fields logrus.Fields) {
	entry := l.withContext(c)
	for k, v := range fields {
		entry = entry.WithField(k, v)
	}
	entry.Error(message)
}

func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.Logger.WithFields(fields)
}
