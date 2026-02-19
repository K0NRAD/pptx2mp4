package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ErrorHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"path":  c.Request.URL.Path,
			}).Error("request-fehler")

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Interner Serverfehler",
				"message": "Ein unerwarteter Fehler ist aufgetreten",
			})
		}
	}
}

func Recovery(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.WithFields(logrus.Fields{
					"error": err,
					"path":  c.Request.URL.Path,
				}).Error("panic aufgetreten")

				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Kritischer Serverfehler",
					"message": "Ein kritischer Fehler ist aufgetreten",
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
